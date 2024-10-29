import socket
import sys
import re

class InputValidator:
    @staticmethod
    def validate_numbers(numbers_str):
        """Validate that input is a sequence of numbers."""
        if not numbers_str.strip():
            return False, "Input cannot be empty"
        
        numbers = numbers_str.strip().split()
        for num in numbers:
            if not num.lstrip('-').isdigit():
                return False, f"Invalid number: {num}"
        return True, numbers_str

    @staticmethod
    def validate_string(string):
        """Validate that input is a non-empty string."""
        if not string.strip():
            return False, "Input cannot be empty"
        return True, string

    @staticmethod
    def validate_sorted_string(string):
        """Validate that input is a sorted sequence."""
        if not string.strip():
            return False, "Input cannot be empty"
        
        elements = string.strip().split()
        if len(elements) < 2:
            return True, string
        
        for i in range(1, len(elements)):
            if elements[i] < elements[i-1]:
                return False, "Sequence is not sorted"
        return True, string

    @staticmethod
    def validate_single_number(number_str):
        """Validate that input is a single positive number."""
        if not number_str.strip():
            return False, "Input cannot be empty"
        
        if not number_str.strip().isdigit():
            return False, "Input must be a positive number"
        return True, number_str

    @staticmethod
    def validate_single_char(char):
        """Validate that input is a single character."""
        char = char.strip()
        if len(char) != 1:
            return False, "Input must be a single character"
        return True, char

    @staticmethod
    def validate_position_and_length(pos_str, len_str):
        """Validate position and length are positive numbers."""
        if not pos_str.strip().isdigit():
            return False, "Position must be a positive number"
        if not len_str.strip().isdigit():
            return False, "Length must be a positive number"
        return True, (pos_str, len_str)

class TCPClient:
    def __init__(self, host, port):
        self.host = host
        self.port = int(port)
        self.validator = InputValidator()
        
    def connect_and_send(self, problem_num, data):
        """Connect to server and send data."""
        try:
            with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as sock:
                sock.connect((self.host, self.port))
                
                # Send problem number first
                sock.sendall(f"{problem_num}\n".encode())
                
                # Send the actual data
                sock.sendall(data.encode())
                
                # Get response
                response = sock.recv(1024).decode()
                return response
        except ConnectionRefusedError:
            return "Error: Could not connect to server"
        except Exception as e:
            return f"Error: {str(e)}"

    def get_validated_input(self, problem_num):
        """Get and validate input based on problem number."""
        try:
            if problem_num == "1":  # Sum of numbers
                print("Enter numbers separated by spaces:")
                numbers = input().strip()
                valid, result = InputValidator.validate_numbers(numbers)
                if not valid:
                    raise ValueError(result)
                return result + "\n"

            elif problem_num == "2":  # Count spaces
                print("Enter a string:")
                text = input().strip()
                valid, result = InputValidator.validate_string(text)
                if not valid:
                    raise ValueError(result)
                return result + "\n"

            elif problem_num == "3":  # Reverse string
                print("Enter a string to reverse:")
                text = input().strip()
                valid, result = InputValidator.validate_string(text)
                if not valid:
                    raise ValueError(result)
                return result + "\n"

            elif problem_num == "4":  # Merge sorted strings
                print("Enter first sorted sequence (elements separated by spaces):")
                seq1 = input().strip()
                valid1, result1 = InputValidator.validate_sorted_string(seq1)
                if not valid1:
                    raise ValueError(f"First sequence: {result1}")
                
                print("Enter second sorted sequence (elements separated by spaces):")
                seq2 = input().strip()
                valid2, result2 = InputValidator.validate_sorted_string(seq2)
                if not valid2:
                    raise ValueError(f"Second sequence: {result2}")
                
                return f"{seq1}\n{seq2}\n"

            elif problem_num == "5":  # Divisors
                print("Enter a positive number:")
                num = input().strip()
                valid, result = InputValidator.validate_single_number(num)
                if not valid:
                    raise ValueError(result)
                return result + "\n"

            elif problem_num == "6":  # Find character positions
                print("Enter a string:")
                text = input().strip()
                valid1, result1 = InputValidator.validate_string(text)
                if not valid1:
                    raise ValueError(result1)
                
                print("Enter a character to search for:")
                char = input().strip()
                valid2, result2 = InputValidator.validate_single_char(char)
                if not valid2:
                    raise ValueError(result2)
                
                return f"{result1}\n{result2}\n"

            elif problem_num == "7":  # Substring
                print("Enter a string:")
                text = input().strip()
                valid1, result1 = InputValidator.validate_string(text)
                if not valid1:
                    raise ValueError(result1)
                
                print("Enter starting position:")
                pos = input().strip()
                print("Enter length:")
                length = input().strip()
                valid2, (pos, length) = InputValidator.validate_position_and_length(pos, length)
                if not valid2:
                    raise ValueError("Invalid position or length")
                
                return f"{result1}\n{pos}\n{length}\n"

            elif problem_num == "8":  # Common numbers
                print("Enter first sequence of numbers (separated by spaces):")
                seq1 = input().strip()
                valid1, result1 = InputValidator.validate_numbers(seq1)
                if not valid1:
                    raise ValueError(f"First sequence: {result1}")
                
                print("Enter second sequence of numbers (separated by spaces):")
                seq2 = input().strip()
                valid2, result2 = InputValidator.validate_numbers(seq2)
                if not valid2:
                    raise ValueError(f"Second sequence: {result2}")
                
                return f"{seq1}\n{seq2}\n"

            elif problem_num == "9":  # Numbers in first but not in second
                print("Enter first sequence of numbers (separated by spaces):")
                seq1 = input().strip()
                valid1, result1 = InputValidator.validate_numbers(seq1)
                if not valid1:
                    raise ValueError(f"First sequence: {result1}")
                
                print("Enter second sequence of numbers (separated by spaces):")
                seq2 = input().strip()
                valid2, result2 = InputValidator.validate_numbers(seq2)
                if not valid2:
                    raise ValueError(f"Second sequence: {result2}")
                
                return f"{seq1}\n{seq2}\n"

            elif problem_num == "10":  # Common characters at same positions
                print("Enter first string:")
                str1 = input().strip()
                valid1, result1 = InputValidator.validate_string(str1)
                if not valid1:
                    raise ValueError(f"First string: {result1}")
                
                print("Enter second string:")
                str2 = input().strip()
                valid2, result2 = InputValidator.validate_string(str2)
                if not valid2:
                    raise ValueError(f"Second string: {result2}")
                
                return f"{result1}\n{result2}\n"

            else:
                raise ValueError("Invalid problem number")

        except ValueError as e:
            print(f"Validation error: {str(e)}")
            return None
        except Exception as e:
            print(f"Error: {str(e)}")
            return None

def main():
    if len(sys.argv) != 4:
        print("Usage: python client.py <host> <port> <problem_number>")
        sys.exit(1)

    host = sys.argv[1]
    port = sys.argv[2]
    problem_num = sys.argv[3]

    if not problem_num.isdigit() or int(problem_num) < 1 or int(problem_num) > 10:
        print("Problem number must be between 1 and 10")
        sys.exit(1)

    client = TCPClient(host, port)
    
    # Get and validate input
    data = client.get_validated_input(problem_num)
    if data is None:
        sys.exit(1)

    # Send data and get response
    response = client.connect_and_send(problem_num, data)
    print("Server response:", response)

if __name__ == "__main__":
    main()