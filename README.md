# Fuzzy Finder

Fuzzy Finder is a command-line tool to search for text within files in a specified directory. It supports both exact and fuzzy matching of search terms, highlighting the matches in the output.

## Features

- Search for exact and fuzzy matches of a given term.
- Highlight matching words in the output.
- Color-coded output for better readability.

## Installation

### Prerequisites

- [Go](https://golang.org/doc/install) (1.16 or later)

### Steps

1. Clone the repository:
    ```sh
    git clone https://github.com/guilhermepelayo/fuzzy.git
    ```

2. Navigate to the project directory:
    ```sh
    cd fuzzy
    ```

3. Build the executable:
    ```sh
    go build -o fuzzy cmd/main.go
    ```

4. Move the executable to a directory in your PATH:
    ```sh
    sudo mv fuzzy /usr/local/bin/
    ```

### Alternative: Adding the Executable to PATH Temporarily

If you prefer not to move the executable, you can add the current directory to your PATH for the session:
    ```sh
    export PATH=$PATH:$(pwd)
    ```

To make this change permanent, add the above line to your `~/.bashrc`, `~/.zshrc`, or other relevant shell configuration file.

## Usage

### Command-Line Arguments

- `search-term`: The term to search for.
- `directory`: The directory to search in.
- `-exact`: (Optional) Use this flag to perform an exact match search.

### Examples

1. Basic usage:
    ```sh
    fuzzy <search-term> </path/to/directory>
    ```

2. Exact match usage:
    ```sh
    fuzzy -exact <search-term> </path/to/directory>
    ```

3. Example:
    ```sh
    fuzzy matx cmd/
    ```

4. Exact match example:
    ```sh
    fuzzy -exact matx cmd/
    ```

### Output

The tool outputs the file path, line number, and the matching line, with the matching word highlighted in yellow, the file path in red, and the line number in green.

## How It Works

- The tool walks through the specified directory and processes each file.
- It checks if the file is a text file.
- It searches for the search term in each line of the file.
- It highlights exact matches and fuzzy matches (words similar to the search term based on Levenshtein distance).