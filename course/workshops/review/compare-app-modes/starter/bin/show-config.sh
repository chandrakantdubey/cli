#!/usr/bin/env bash
printf 'mode=%s\nlog=%s\n' "${APP_MODE:-local}" "${LOG_LEVEL:-info}"
