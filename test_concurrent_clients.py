import argparse
import subprocess
import threading
import time
import random
from concurrent.futures import ThreadPoolExecutor, as_completed
import sys
from datetime import datetime

class ConcurrencyTester:
    def __init__(self, host, port, num_clients, test_duration=30):
        self.host = host
        self.port = port
        self.num_clients = num_clients
        self.test_duration = test_duration
        self.results = []
        self.lock = threading.Lock()
        self.start_time = None

    def run_client(self, client_id):
        """Executa un client individual"""
        start_time = time.time()
        result = {
            'client_id': client_id,
            'success': False,
            'duration': 0,
            'elapsed': 0,
            'error': None,
            'output': None,
            'problem': None
        }

        try:
            problem = random.randint(1, 10)
            result['problem'] = problem
            input_data = self.generate_input(problem)
            
            cmd = ["python3", "client/client.py", self.host, str(self.port), str(problem)]
            
            process = subprocess.Popen(
                cmd,
                stdin=subprocess.PIPE,
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE,
                text=True
            )
            
            stdout, stderr = process.communicate(input=input_data, timeout=10)
            
            result['duration'] = time.time() - start_time
            result['elapsed'] = time.time() - self.start_time
            result['success'] = process.returncode == 0
            result['output'] = stdout.strip()
            result['error'] = stderr.strip() if stderr else None
            
            with self.lock:
                self.results.append(result)
                self.print_progress(client_id, result['success'], result['duration'], problem)
                
        except subprocess.TimeoutExpired:
            result['error'] = "Client timeout"
            result['duration'] = time.time() - start_time
            result['elapsed'] = time.time() - self.start_time
            with self.lock:
                self.results.append(result)
                self.print_progress(client_id, False, result['duration'], result.get('problem', 'N/A'))

        except Exception as e:
            result['error'] = str(e)
            result['duration'] = time.time() - start_time
            result['elapsed'] = time.time() - self.start_time
            with self.lock:
                self.results.append(result)
                self.print_progress(client_id, False, result['duration'], result.get('problem', 'N/A'))

    def generate_input(self, problem):
        """Genereaza input random pentru fiecare tip de problema"""
        if problem == 1:  # Sum of numbers
            nums = [random.randint(1, 100) for _ in range(5)]
            return " ".join(map(str, nums)) + "\n"
            
        elif problem == 2:  # Count spaces
            words = ["hello", "world", "this", "is", "a", "test"]
            return " ".join(random.choices(words, k=4)) + "\n"
            
        elif problem == 3:  # Reverse string
            words = ["python", "golang", "testing", "concurrent", "server"]
            return random.choice(words) + "\n"
            
        elif problem == 4:  # Merge sorted strings
            seq1 = sorted([random.randint(1, 100) for _ in range(3)])
            seq2 = sorted([random.randint(1, 100) for _ in range(3)])
            return f"{' '.join(map(str, seq1))}\n{' '.join(map(str, seq2))}\n"
            
        elif problem == 5:  # Divisors
            return f"{random.randint(1, 100)}\n"
            
        elif problem == 6:  # Find character positions
            text = "hello world testing"
            char = random.choice("abcdefghijklmnopqrstuvwxyz")
            return f"{text}\n{char}\n"
            
        elif problem == 7:  # Substring
            text = "concurrent server test"
            pos = random.randint(0, len(text)-1)
            length = random.randint(1, 5)
            return f"{text}\n{pos}\n{length}\n"
            
        elif problem == 8:  # Common numbers
            arr1 = [random.randint(1, 100) for _ in range(5)]
            arr2 = [random.randint(1, 100) for _ in range(5)]
            return f"{' '.join(map(str, arr1))}\n{' '.join(map(str, arr2))}\n"
            
        elif problem == 9:  # Numbers in first but not in second
            arr1 = [random.randint(1, 100) for _ in range(5)]
            arr2 = [random.randint(1, 100) for _ in range(5)]
            return f"{' '.join(map(str, arr1))}\n{' '.join(map(str, arr2))}\n"
            
        else:  # Most common character at same positions
            str1 = "hello world"
            str2 = "happy world"
            return f"{str1}\n{str2}\n"

    def print_progress(self, client_id, success, duration, problem):
        """Afiseaza progres live"""
        status = "✓" if success else "✗"
        color = "\033[92m" if success else "\033[91m"
        reset = "\033[0m"
        print(f"{color}Client {client_id:3d} {status} Problem {problem} ({duration:.2f}s){reset}")

    def print_summary(self):
        """Afiseaza sumarul testului"""
        if not self.results:
            print("\nNo results to display - all clients failed to connect")
            return

        successful = sum(1 for r in self.results if r.get('success', False))
        failed = len(self.results) - successful
        
        # Calculam doar pentru rezultatele care au elapsed si duration
        valid_results = [r for r in self.results if 'elapsed' in r and 'duration' in r]
        if valid_results:
            total_duration = max(r['elapsed'] for r in valid_results)
            avg_duration = sum(r['duration'] for r in valid_results) / len(valid_results)
        else:
            total_duration = 0
            avg_duration = 0

        print("\n" + "="*50)
        print(f"Concurrency Test Summary")
        print("="*50)
        print(f"Total Clients:     {self.num_clients}")
        print(f"Successful:        {successful}")
        print(f"Failed:           {failed}")
        print(f"Success Rate:     {(successful/self.num_clients)*100:.1f}%")
        print(f"Total Duration:   {total_duration:.2f}s")
        print(f"Average Duration: {avg_duration:.2f}s")
        if total_duration > 0:
            print(f"Requests/second:  {self.num_clients/total_duration:.1f}")
        print("="*50)

        # Afiseaza primele 5 erori (daca exista)
        errors = [r for r in self.results if r.get('error')]
        if errors:
            print("\nSample Errors:")
            for e in errors[:5]:
                print(f"Client {e['client_id']}: {e['error']}")

    def run_test(self):
        """Ruleaza testul complet"""
        print(f"Starting concurrency test with {self.num_clients} clients...")
        print(f"Target server: {self.host}:{self.port}\n")
        
        self.start_time = time.time()
        
        with ThreadPoolExecutor(max_workers=self.num_clients) as executor:
            futures = [
                executor.submit(self.run_client, i) 
                for i in range(self.num_clients)
            ]
            
            for future in as_completed(futures):
                try:
                    future.result()
                except Exception as e:
                    print(f"Error in client thread: {e}")

        self.print_summary()

def main():
    parser = argparse.ArgumentParser(description='Test TCP Server Concurrency')
    parser.add_argument('--host', default='localhost', help='Server host')
    parser.add_argument('--port', type=int, default=8080, help='Server port')
    parser.add_argument('--clients', type=int, default=10, help='Number of concurrent clients')
    parser.add_argument('--duration', type=int, default=30, help='Test duration in seconds')
    
    args = parser.parse_args()
    
    tester = ConcurrencyTester(
        args.host,
        args.port,
        args.clients,
        args.duration
    )
    
    tester.run_test()

if __name__ == "__main__":
    main()