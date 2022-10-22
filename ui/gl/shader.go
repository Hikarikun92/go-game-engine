package gl

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"strings"
)

// Vertex shader using the projection matrix (defining the screen size and orientation), the model's vertices and the
// model's attributes (position, size etc.). The X and Y coordinates of each vertex are the position of the vertex, and
// the last 2 coordinates (Z and W) are the coordinates of the associated texture.
var vertexShader = `
#version 330 core
layout (location = 0) in vec4 vertexData;

out vec2 TexCoord;

uniform mat4 model;
uniform mat4 projection;

void main()
{
	gl_Position = projection * model * vec4(vertexData.xy, 0.0, 1.0);
	TexCoord = vertexData.zw;
}
` + "\x00"

// Fragment shader that will simply retrieve the texture's color at the specified point.
var fragmentShader = `
#version 330

uniform sampler2D tex;

in vec2 TexCoord;

out vec4 outputColor;

void main() {
    outputColor = texture(tex, TexCoord);
}
` + "\x00"

func newShaderProgram() (uint32, error) {
	vertexShader, err := compileShader(vertexShader, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShader, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		logValue := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(logValue))

		return 0, fmt.Errorf("failed to link program: %v", logValue)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	cSources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, cSources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		logValue := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(logValue))

		return 0, fmt.Errorf("failed to compile %v: %v", source, logValue)
	}

	return shader, nil
}
