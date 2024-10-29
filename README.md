# TCP Problem Solver

O implementare client-server TCP. Include atat serverul (Go) cat si clientul (Python).

## Structura

```
├── server/
│   └── main.go    # Serverul TCP implementat in Go
└── client/
    └── client.py  # Clientul TCP implementat in Python
```

## Cerinte

- Go 1.16 sau mai nou (pentru server)
- Python 3.7 sau mai nou (pentru client)

## Functionalitati

Serverul poate rezolva urmatoarele probleme:

1. Suma numerelor
2. Numararea spatiilor
3. Inversarea unui sir
4. Interclasarea sirurilor sortate
5. Divizorii unui numar
6. Gasirea pozitiilor unui caracter
7. Substring
8. Numere comune
9. Numere prezente in primul array dar nu in al doilea
10. Cel mai frecvent caracter de pe aceeasi pozitie

## Utilizare

### Server

```bash
# Compilare
go build -o server main.go

# Rulare (specificand host si port optional)
./server -host=0.0.0.0 -port=8080
```

### Client

```bash
# Rulare client (specificand host, port si numarul problemei)
python client.py <host> <port> <problem_number>

# Exemplu:
python client.py localhost 8080 1
```

## Exemple

```python
# Problema 1: Suma numerelor
python client.py localhost 8080 1
# Input: 1 2 3 4 5
# Output: Sum: 15

# Problema 3: Inversare sir
python client.py localhost 8080 3
# Input: Hello World
# Output: dlroW olleH
```
