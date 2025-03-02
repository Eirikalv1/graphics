struct VertexInput {
    @location(0) position: vec2<f32>,
    @location(1) color: vec2<f32>,    
}

struct VertexOutput {
    @builtin(position) clip_position: vec4<f32>,
    @location(0) color: vec2<f32>
}

@vertex
fn vs_main(model: VertexInput) -> VertexOutput {
    var out: VertexOutput;
    out.clip_position = vec4<f32>(model.position, 0.0, 1.0);
    out.color = model.color;
    return out;
}

@fragment
fn fs_main(in: VertexOutput) -> @location(0) vec4<f32> {
    let coord = vec2<f32>(in.color) * 2.0 - 1.0; // -1 -> 1
    let color = per_pixel(coord);
    return color;
}

struct Ray {
    origin: vec3<f32>,
    direction: vec3<f32>
}

fn new_ray(coord: vec2<f32>) -> Ray {
    var ray: Ray;

    ray.direction = vec3<f32>(coord, -1.0);
    ray.origin = vec3<f32>(0.0, 0.0, 1.0);
    return ray;
}

struct Sphere {
    albedo: vec3<f32>,
    position: vec3<f32>,
    radius: f32,
}

fn new_sphere() -> Sphere {
    var sphere: Sphere;
    sphere.albedo = vec3<f32>(1.0, 0.0, 1.0);
    sphere.position = vec3<f32>(0.0);
    sphere.radius = 0.5;
    return sphere;
}

fn per_pixel(coord: vec2<f32>) -> vec4<f32> {
    let ray = new_ray(coord);
    let sphere = new_sphere();

    let origin = ray.origin - sphere.position;

    let a = dot(ray.direction, ray.direction);
    let b = 2.0 * dot(origin, ray.direction);
    let c = dot(origin, origin) - pow(sphere.radius, 2.0);

    let discriminant = b * b - 4.0 * a * c;
    if discriminant < 0.0 {
        return vec4<f32>(0.0, 0.0, 0.0, 1.0);
    }

    let closestT = (-b - sqrt(discriminant)) / (2.0 * a);
    let t0 = (-b + sqrt(discriminant)) / (2.0 * a);

    let hit_point = origin + ray.direction * closestT;
    let hit_normal = normalize(hit_point);

    let light_dir = normalize(vec3<f32>(0.0, 0.0, -1.0));
    let light_intensity = max(dot(hit_normal, -light_dir), 0.0);

    var sphere_color = sphere.albedo;
    sphere_color *= light_intensity;

    return vec4<f32>(sphere_color, 1.0);
}