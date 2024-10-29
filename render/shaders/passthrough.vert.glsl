#version 460 core
layout (location = 0) in vec3 aPos;

uniform vec3 res;

void main()
{
	// Convert from world space to the range -1 to 1
	vec3 uv = ((aPos / res) - 0.5) * 2;

    gl_Position = vec4(uv.x, uv.y, aPos.z, 1.0);
}
