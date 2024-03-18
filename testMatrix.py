import requests
import time
import argparse
import concurrent.futures
import csv
import random


# NOTE
# in order to multiply two matrices, the number of columns in matrix A
# must be equal to the number of rows in column B.


def make_request(url, rows, cols):
    response = requests.get(f"{url}?rows={rows}&cols={cols}")
    return response.json()

def make_parallel_request(url, a, b, start, end):
    data = {
        "a": a,
        "b": b,
        "start": start,
        "end": end
    }
    response = requests.post(url, json=data)

    return response.json()

def measure_time(url, rows, cols):
    start = time.time()
    result = make_request(url, rows, cols)
    elapsed = time.time() - start
    return result, elapsed

def write_to_file(filename, data, csv_filename, csv_row):
    with open(filename, 'a') as f:
        f.write(data)
    
    with open(csv_filename, 'a', newline='') as f:
        writer = csv.writer(f)
        writer.writerow(csv_row)

def main():
    parser = argparse.ArgumentParser(description='Performance test for URLs.')
    parser.add_argument('--rows', type=int, default=10, help='Number of rows in matrix.')
    parser.add_argument('--cols', type=int, default=10, help='Number of cols in matrix.')
    parser.add_argument('--processes', type=int, default=1, help='Number of processes for the multifunction test.')

    args = parser.parse_args()

    matrix_sizes = [2, 4, 8, 16, 32, 64, 128, 256, 512]
    numProcs = [1, 2, 4, 8]


    rows = args.rows
    cols = args.cols
    processes = args.processes

    sequentialUrls = {
        "Azure Sequential": "https://sequential-function.azurewebsites.net/api/sequential-processing",
        "Azure Threaded": "https://threaded-function.azurewebsites.net/api/in-function-parallelism",
        "OpenFaaS Sequential": "http://20.26.236.208:8080/function/sequential-matrix",
        "OpenFaaS Threaded": "http://20.26.236.208:8080/function/threaded-matrix",
    }

    for size in matrix_sizes:
        rows = size
        cols = size

        for name, url in sequentialUrls.items():
            # Test by passing all times at once
            result = make_request(url, rows, cols)
            if 'numProcs' in result:
                write_to_file("results.txt", f"{name}, Matrix Size: {result['matrix size']}, Elapsed time: {result['execTime']}, Threads: {result['numProcs']}\n", 
                            "results.csv", [name, result['matrix size'], result['execTime'], result['numProcs']])
            else:
                write_to_file("results.txt", f"{name}, Matrix Size: {result['matrix size']}, Elapsed time: {result['execTime']}\n", 
                            "results.csv", [name, result['matrix size'], result['execTime']])
        
        multifunctionUrls = {
            "Azure Multi-function": "https://multi-function-parallelism.azurewebsites.net/api/multi-function-parallelism?",
            "OpenFaaS Multi-function": "http://20.26.236.208:8080/function/multi-fuction-matrix",
        }

            # Initialize matrix a
        a = [[random.randint(10, 100) for _ in range(cols)] for _ in range(rows)]

        # Initialize matrix b
        b = [[random.randint(10, 100) for _ in range(cols)] for _ in range(cols)]

        # Initialize result matrix
        result = [[0 for _ in range(cols)] for _ in range(rows)]

        for proc in numProcs:
            processes = proc
            with concurrent.futures.ThreadPoolExecutor() as executor:
                for name, url in multifunctionUrls.items():
                    result = [[0 for _ in range(cols)] for _ in range(rows)]
                    startTime = time.time()
                    rows_per_process = rows // processes
                    cols_per_process = cols // processes
                    futures = []
                    for i in range(processes):
                        start = i * rows_per_process
                        end = (start + rows_per_process) if i < processes - 1 else rows - 1
                        future = executor.submit(make_parallel_request, url, a, b, start, end)
                        futures.append(future)
                    totalMatrix = sum((f.result()['result'] for f in concurrent.futures.as_completed(futures)), [])
                    totalMatrixSize = rows * cols
                    elapsed = time.time() - startTime
                    elapsedSec = "{:.8f}".format(elapsed)
                    write_to_file("results.txt", f"{name}, Matrix Size: {totalMatrixSize}, Elapsed time: {elapsedSec}, Processes: {processes}\n", 
                            "results.csv", [name, totalMatrixSize, elapsedSec, processes])

if __name__ == "__main__":
    main()