FROM python:3

ENV PIP_DISABLE_PIP_VERSION_CHECK=on

RUN pip install poetry

WORKDIR /usr/src/app

RUN mkdir ./crawler
COPY ./crawler ./crawler
COPY ./poetry.lock ./pyproject.toml ./

RUN poetry config virtualenvs.create false
RUN poetry install --no-interaction

CMD ["poetry", "run", "python", "./crawler"]