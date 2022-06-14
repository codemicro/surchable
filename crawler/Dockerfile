FROM python:3-alpine

ENV PIP_DISABLE_PIP_VERSION_CHECK=on

RUN pip install poetry

WORKDIR /usr/src/app

COPY ../poetry.lock ../pyproject.toml ./

RUN poetry config virtualenvs.create false
RUN poetry install --no-interaction

COPY ./crawler ./

CMD ["poetry", "run", "python", "./crawler"]