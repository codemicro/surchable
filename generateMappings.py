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


def generate_golang_url_mapping(mapping: Dict[str, str], package_name: str) -> str:
    o = f"""package {package_name}\n\nconst (\n"""

    for key in mapping:
        o += f'\t{key} = "{mapping[key]}"\n'

    o += ")\n"

    return o


def generate_python_url_mapping(mapping: Dict[str, str]):
    o = ""
    for key in mapping:
        o += f'{_camel_to_snake(key).upper()} = "{mapping[key]}"\n'
    return o


def parse_configuration(raw: str) -> Dict[str, Tuple[str, Optional[str]]]:
    mappings: Dict[str, Tuple[str, Optional[str]]] = None

    for i, line in enumerate(_fields(raw)):
        if line == "":
            continue

        fields = _fields(line)

        assert len(fields) >= 3, str(i+1)

        if fields[0] == "must":
            mappings[fields[1]] = (fields[2], None)
        elif fields[1] == "default":
            mappings[fields[2]] = (fields[3], fields[1])

    return mappings

def generate_golang_configuration(mapping: Dict[str, Tuple[str, Optional[str]]], package_name: str) -> str:
    raise NotImplementedError("generate_golang_configuration")

def generate_python_configuration(mapping: Dict[str, Tuple[str, Optional[str]]]) -> str:
    raise NotImplementedError("generate_python_configuration")


def run_url_mappings():
    mappings: Dict[str, str] = {}
    with open("mappings/urls") as f:
        mappings = parse_url_mappings(f.read())

    with open("internal/urls/urls.go", "w") as f:
        f.write(generate_golang_url_mapping(mappings, "urls"))

    with open("crawler/urls.py", "w") as f:
        f.write(generate_python_url_mapping(mappings))


def run_configuration():
    mappings: Dict[str, Tuple[str, Optional[str]]] = {}
    with open("mappings/configuration") as f:
        mappings = parse_configuration(f.read())

    with open("internal/config/envVars.go", "w") as f:
        f.write(generate_golang_configuration(mappings, "urls"))

    with open("crawler/envVars.py", "w") as f: # TODO: Change this to replace part of the file instead.
        f.write(generate_python_url_mapping(mappings))


# TODO: Reimplement this with Cog preprocessor??


if __name__ == "__main__":
    run_url_mappings()
