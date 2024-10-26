package render

import (
	"fmt"
	"os"

	"github.com/go-gl/gl/v4.5-core/gl"
)

func LoadShaderFromFile(path string, shaderType uint32) (uint32, error) {
	// Load from files into strings
	_shaderSource, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}
	shaderSource := string(_shaderSource)

	fmt.Printf(shaderSource)

	return compileShader(shaderSource, shaderType)
}

func LoadShaderProgram(vertShaderPath, fragShaderPath string) (uint32, error) {
	vertShader, err := LoadShaderFromFile(vertShaderPath, gl.VERTEX_SHADER)
	defer gl.DeleteShader(vertShader)
	if err != nil {
		return 0, err
	}

	fragShader, err := LoadShaderFromFile(fragShaderPath, gl.FRAGMENT_SHADER)
	defer gl.DeleteShader(fragShader)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertShader)
	gl.AttachShader(program, fragShader)

	gl.LinkProgram(program)
	defer gl.DetachShader(program, vertShader)
	defer gl.DetachShader(program, fragShader)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var length int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &length)

		log := string(make([]byte, length))
		gl.GetProgramInfoLog(program, length, nil, gl.Str(log))
		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	// Create shader object
	shader := gl.CreateShader(shaderType)

	// Set shader source
	csource, free := gl.Strs(source)
	defer free()

	gl.ShaderSource(shader, 1, csource, nil)

	// Compile shader
	gl.CompileShader(shader)

	// Check for compilation errors
	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var length int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &length)

		log := string(make([]byte, length))
		gl.GetShaderInfoLog(shader, length, nil, gl.Str(log))
		return 0, fmt.Errorf("failed to compile shader: %v", log)
	}

	return shader, nil
}
