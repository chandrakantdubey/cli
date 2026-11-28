#!/usr/bin/env python3
"""Validate authored course structure without running lesson commands."""
import json
import re
import sys
from pathlib import Path


def validate(root):
    errors = []
    catalog = json.loads((root / 'course/catalog.json').read_text())
    ids = [ident for group in catalog['workshops'] for ident in group['lessons']]
    if len(ids) != len(set(ids)):
        errors.append('Catalog has duplicate IDs')
    for ident in ids:
        if not re.fullmatch(r'[a-z0-9-]+/[a-z0-9-]+', ident):
            errors.append(f'Invalid ID: {ident}')
            continue
        directory = root / 'course/workshops' / ident
        try:
            lesson = json.loads((directory / 'lesson.json').read_text())
            body = (directory / 'lesson.md').read_text()
        except (OSError, ValueError) as error:
            errors.append(f'{ident}: {error}')
            continue
        if lesson['id'] != ident:
            errors.append(f'{ident}: manifest identity differs')
        if sum(item['points'] for item in lesson['checks']) != 100:
            errors.append(f'{ident}: rubric must total 100')
        required = ['Commands and operations']
        if not lesson.get('reviewOf'):
            required += ['Before you practice', 'Common mistakes and recovery', 'Use it in development']
        for heading in required:
            if f'## {heading}\n' not in body:
                errors.append(f'{ident}: missing teaching section {heading}')
        if len(lesson.get('hints', [])) != 3 or any(not h.strip() for h in lesson['hints']):
            errors.append(f'{ident}: three progressive hints are required')
        for frame in lesson.get('demo', []):
            if any(not frame.get(key, '').strip() for key in ('command', 'visual', 'explanation')):
                errors.append(f'{ident}: incomplete worked example')
        base = lesson.get('reviewOf')
        if base and (base not in ids or base.startswith('review/')):
            errors.append(f'{ident}: invalid original lesson {base}')
        for target in re.findall(r'\]\(/workshop/([a-z0-9/-]+)', body):
            if target not in ids:
                errors.append(f'{ident}: broken workshop link {target}')
    return errors


if __name__ == '__main__':
    root = Path(sys.argv[1]).resolve() if len(sys.argv) > 1 else Path(__file__).resolve().parents[1]
    errors = validate(root)
    if errors:
        print('\n'.join(errors), file=sys.stderr)
        sys.exit(1)
    print('Course teaching sections, identities, links, hints, examples, and point totals verified.')
