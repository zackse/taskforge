"""Benchmark parser code"""

import cProfile

from taskforge.ql.parser import Parser


def benchmark_parser():
    """Benchmark the performance of various queries."""
    queries = [
        'milk and cookies',
        'milk -and cookies',
        'completed = false',
        '(priority > 5 and title ^ \'take out the trash\') or '
        '(context = "work" and (priority >= 2 or ("my little pony")))',
    ]

    for query in queries:
        profiler = cProfile.Profile()
        profiler.enable()
        parser = Parser(query)
        parser.parse()
        profiler.disable()
        print(f'{query}:')
        profiler.print_stats()


if __name__ == '__main__':
    benchmark_parser()
