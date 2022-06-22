import re
from typing import *

_camel_to_snake_pattern = re.compile(r"(?<!^)(?=[A-Z])")


def _camel_to_snake(inp: str) -> str:
    if inp.upper() == inp:
        return inp
    return _camel_to_snake_pattern.sub("_", inp)


def _fields(inp: str) -> List[str]:
    return list(filter(lambda x: x != "", inp.split(" ")))


def parse_url_mappings(raw: str) -> Dict[str, str]:
    mappings: Dict[str, str] = {}

    for i, line in enumerate(raw.strip().splitlines()):
        if line == "":
            continue
        fields = _fields(line)
        assert (
            len(fields) == 2
        ), f"Line {i+1} has an invalid number of fields, {len(fields)}. It must have 2."

        mappings[fields[0]] = fields[1]

    return mappings


def generate_golang_url_mapping(mapping: Dict[str, str]) -> str:
    o = "\n\n// The below was generated. Do not edit.\n// Modify mappings/urls instead.\n\nconst (\n"

    for key in mapping:
        o += f'\t{key} = "{mapping[key]}"\n'

    o += ")\n"

    return o


def generate_python_url_mapping(mapping: Dict[str, str]):
    o = "# The below was generated. Do not edit.\n# Modify mappings/urls instead.\n"
    for key in mapping:
        o += f'{_camel_to_snake(key).upper()} = "{mapping[key]}"\n'
    return o


def parse_configuration(raw: str) -> Dict[str, Tuple[str, Optional[str]]]:
    mappings: Dict[str, Tuple[str, Optional[str]]] = {}

    for i, line in enumerate(raw.splitlines()):
        if line == "":
            continue

        fields = _fields(line)

        assert len(fields) >= 3, str(i+1)

        if fields[0] == "must":
            mappings[fields[1]] = (fields[2], None)
        elif fields[0] == "default":
            mappings[fields[2]] = (fields[3], fields[1])

    return mappings

def generate_golang_configuration(mapping: Dict[str, Tuple[str, Optional[str]]]) -> str:
    highLevelTypes: Dict[str, Tuple[str, str, Optional[str]]] = {}
    for key in mapping:
        kc = key.count(".")
        sp = key.split(".")
        if kc == 1:
            x = highLevelTypes.get(sp[0], [])
            if type(x) != list:
                raise ValueError(f"duplicated key {sp}")
            x.append((sp[1], mapping[key][0], mapping[key][1]))
            highLevelTypes[sp[0]] = x
        elif kc == 0:
            highLevelTypes[key] = (key, mapping[key][0], mapping[key][1])
        else:
            raise NotImplementedError("subtypes more than one level deep")

    o = "\n\n// The below was generated. Do not edit.\n// Modify mappings/configuration instead.\n\n"

    def formFunctionCall(x: Tuple[str, str, Optional[str]]) -> str:
        return ('requireEnvVar' if x[2] is None else 'envVarDefault') + "(\"" + x[1] + "\"" + (', \"' + x[2] + '\"' if x[2] is not None else '') + ")"

    for key in highLevelTypes:
        val = highLevelTypes[key]
        if type(val) == tuple:
            o += f"var {key} = {formFunctionCall(val)})\n"
        if type(val) == list:
            o += f"var {key} = struct{{\n"
            for var in val:
                o += f"{var[0]} string\n"
            o += "}{\n"
            for var in val:
                o += f"{var[0]}: {formFunctionCall(var)},\n"
            o += "}\n"

    return o

def generate_python_configuration(mapping: Dict[str, Tuple[str, Optional[str]]]) -> str:
    o = "# The below was generated. Do not edit.\n# Modify mappings/urls instead.\n"
    for key in mapping:
        x = mapping[key]
        joined_key = ("_".join(map(_camel_to_snake, key.split(".")))).upper()
        o += joined_key + " = " + ('_required_env_var' if x[1] is None else '_env_var_default') + "(\"" + x[0] + "\"" + (', \"' + x[1] + '\"' if x[1] is not None else '') + ")\n"
    return o


def load_raw_url_mappings() -> str:
    with open("mappings/urls") as f:
        return f.read()


def load_raw_configuration() -> str:
    with open("mappings/configuration") as f:
        return f.read()
