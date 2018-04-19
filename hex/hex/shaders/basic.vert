#version 410

layout(location = 0) in vec3 vertex_pos;
layout(location = 1) in vec4 vertexColor;
out vec4 vColor;

void main() {
  vColor = vertexColor;
  gl_Position = vec4(vertex_pos, 1.0);
}
