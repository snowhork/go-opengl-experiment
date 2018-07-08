#version 410

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

layout(location = 0) in vec3 vertex_pos;
layout(location = 1) in vec4 vertexColor;
out vec4 vColor;

void main() {
  vColor = vertexColor;

  gl_Position = projection * camera * model * vec4(vertex_pos, 1);

}
