#include <iostream>
#include <string.h> // For strcpy

void vulnerable_function(char* input) {
    char buffer[10];
    strcpy(buffer, input); // Buffer overflow vulnerability
    std::cout << "Buffer: " << buffer << std::endl;
}

int main(int argc, char** argv) {
    if (argc < 2) {
        std::cout << "Usage: " << argv[0] << " <input_string>" << std::endl;
        return 1;
    }
    vulnerable_function(argv[1]);
    return 0;
}
