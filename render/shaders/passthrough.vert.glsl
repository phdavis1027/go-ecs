#version 460 core
layout (location = 0) in vec3 aPos;

uniform mat4 mvp;

void main()
{
	// Convert from world space to the range -1 to 1
	gl_Position = mvp * vec4(aPos, 1.0);

}
