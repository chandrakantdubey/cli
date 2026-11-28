"""Exercise authoring failures without modifying the published course."""
import importlib.util
import json
import tempfile
import unittest
from pathlib import Path

spec = importlib.util.spec_from_file_location('quality', Path(__file__).with_name('course-quality.py'))
quality = importlib.util.module_from_spec(spec)
spec.loader.exec_module(quality)


class ContentChecks(unittest.TestCase):
    def test_rejects_incomplete_teaching_and_broken_rubric(self):
        with tempfile.TemporaryDirectory() as temporary:
            root = Path(temporary)
            directory = root / 'course/workshops/example/lesson'
            directory.mkdir(parents=True)
            (root / 'course/catalog.json').write_text(json.dumps({'workshops': [{'lessons': ['example/lesson']}]}))
            (directory / 'lesson.json').write_text(json.dumps({'id': 'example/lesson', 'hints': ['one', 'two', 'three'], 'checks': [{'points': 50}], 'demo': []}))
            (directory / 'lesson.md').write_text('## Commands and operations\n\nAn example.\n')
            errors = quality.validate(root)
            self.assertTrue(any('total 100' in error for error in errors))
            self.assertTrue(any('Before you practice' in error for error in errors))


if __name__ == '__main__':
    unittest.main()
