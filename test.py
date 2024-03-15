import requests
import time
import argparse
import concurrent.futures
import csv


def make_request(url, times):
    response = requests.get(f"{url}?times={times}")
    return response.json()

def measure_time(url, times):
    start = time.time()
    result = make_request(url, times)
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
    parser.add_argument('--times', type=int, default=1000, help='Number of times to hit each URL.')
    parser.add_argument('--processes', type=int, default=1, help='Number of processes for the sequential test.')
    args = parser.parse_args()

    times = args.times
    processes = args.processes
    timesPerProc = times // processes


    sequentialUrls = {
        "Azure Sequential": "https://sequential-function.azurewebsites.net/api/sequential-processing",
        "Azure Threaded": "https://threaded-function.azurewebsites.net/api/in-function-parallelism",
        "OpenFaaS Sequential": "http://20.26.236.208:8080/function/sequential-function",
        "OpenFaaS Threaded": "http://20.26.236.208:8080/function/threaded-processing",
    }

    for name, url in sequentialUrls.items():
        # Test by passing all times at once
        result = make_request(url, times)
        write_to_file("results.txt", f"{name},  Times: {times}, Result: {result['result']}, Elapsed time: {result['execTime']}\n", 
                        "results.csv", [name, times, result['execTime'], processes])
        # Test by passing a fraction of times to a defined number of processes
        # for _ in range(processes):
        #     result, elapsed = measure_time(url, times // processes)
        #     write_to_file("results.txt", f"Name: {name}, URL: {url}, Times: {times // processes}, Result: {result}, Elapsed time: {elapsed}\n")

    multifunctionUrls = {
        "Azure Sequential": "https://sequential-function.azurewebsites.net/api/sequential-processing",
        # "Azure Threaded": "https://threaded-function.azurewebsites.net/api/in-function-parallelism",
        "OpenFaaS Sequential": "http://20.26.236.208:8080/function/sequential-function",
        # "OpenFaaS Threaded": "http://20.26.236.208:8080/function/threaded-processing",
    }

    with concurrent.futures.ThreadPoolExecutor() as executor:
        for name, url in multifunctionUrls.items():
            startTime = time.time()
            futures = [executor.submit(make_request, url, timesPerProc) for _ in range(processes)]
            totalScore = sum(f.result()['result'] for f in concurrent.futures.as_completed(futures))
            elapsed = time.time() - startTime
            elapsedSec = "{:.8f}".format(elapsed)
            write_to_file("results.txt", f"URL: {name}, Total Score: {totalScore}, Elapsed time: {elapsedSec}\n", 
                      "results.csv", [name, times, elapsedSec, processes])

if __name__ == "__main__":
    main()