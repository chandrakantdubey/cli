#!/usr/bin/env python3
"""Create a reviewable lesson draft outside the published course."""
import argparse
import json
import re
from pathlib import Path

parser = argparse.ArgumentParser(description=__doc__)
parser.add_argument('id', help='group/slug')
parser.add_argument('--output', type=Path, default=Path('.local/lesson-drafts'))
args = parser.parse_args()
if not re.fullmatch('[a-z0-9-]+/[a-z0-9-]+', args.id):
    parser.error('Use a lowercase group/slug without dots or escaping paths')
dest = args.output / args.id
if dest.exists():
    parser.error(f'Draft already exists: {dest}')
(dest / 'starter').mkdir(parents=True)
(dest / 'starter/message.txt').write_text('Service ready.\n')
lesson = {
    'id': args.id, 'title': 'Save a service note', 'starterRevision': 1, 'gradingRevision': 1,
    'task': 'Read message.txt and save its intact contents as report.txt. Preserve message.txt.',
    'hints': ['Inspect message.txt with cat.', 'Redirect stdout to report.txt.', 'Run cat message.txt > report.txt, then cli-silly check.'],
    'checks': [{'id': 'report', 'kind': 'text', 'path': 'report.txt', 'points': 100,
                'expected': 'Service ready.\n', 'label': 'Report matches the original note', 'preserve': ['message.txt'],
                'feedback': 'Read both files and check the final newline.'}],
    'demo': [{'command': 'cat note.txt > copy.txt', 'visual': 'note.txt -> cat -> copy.txt',
              'explanation': 'Bash opens copy.txt and sends cat stdout to it.'}],
}
(dest / 'lesson.json').write_text(json.dumps(lesson, indent=2) + '\n')
(dest / 'lesson.md').write_text('''Copy a service note while preserving its contents.

## Before you practice

Know how to find your current directory with pwd. Bash and cat are sufficient.

## Commands and operations

cat reads its file argument. The shell operator > creates or truncates the output file and sends stdout to it.

## Common mistakes and recovery

A command that redirects output may print nothing. Read the output file to verify it. Relative output paths start in the current directory.

## Use it in development

Capture a service diagnostic for a colleague while keeping its original input available.
''')
print(f'Draft created at {dest}. Adapt its actual teaching and task, add grader scenarios, then register it in course/catalog.json. See docs/lab-authoring.md.')
