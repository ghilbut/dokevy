# syntax=docker/dockerfile:1.5.1

FROM --platform=$BUILDPLATFORM python:3.11.2-slim AS base

WORKDIR /app

ENV LANG   C.UTF-8
ENV LC_ALL C.UTF-8
ENV PYTHONDONTWRITEBYTECODE 1
ENV PYTHONFAULTHANDLER      1


FROM base AS deps

RUN pip install pipenv

COPY Pipfile .
COPY Pipfile.lock .
RUN pipenv install --system


FROM base AS builder

COPY --from=deps /usr/local/lib/python3.11/site-packages /usr/local/lib/python3.11/site-packages
COPY static ./static
COPY k8single ./k8single
COPY manage.py .

RUN python ./manage.py collectstatic --clear --no-input


FROM base AS release
LABEL maintainer="ghilbut@gmail.com"
LABEL repository="https://github.com/ghilbut/docker-desktop"

EXPOSE 8000

COPY --from=deps /usr/local/lib/python3.11/site-packages /usr/local/lib/python3.11/site-packages
COPY --from=deps /usr/local/bin/uvicorn /usr/local/bin/uvicorn
COPY --from=builder /app/static /app/static
COPY --from=builder /app/k8single /app/k8single
CMD ["uvicorn", "--host", "0.0.0.0", "k8single.asgi:application"]
