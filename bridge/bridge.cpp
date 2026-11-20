#include "bridge.h"
#include <cstdio>
#include <cstdlib>
#include <cstring>
#include <filesystem>
#include <fstream>
#include <memory>
#include <optional>
#include <string>
#include <vector>

#include <slint-interpreter.h>
#include <slint_image.h>
#include <slint_string.h>

using namespace slint::interpreter;

static std::unique_ptr<ComponentCompiler> g_compiler;
static std::optional<slint::ComponentHandle<ComponentInstance>> g_instance;
static std::vector<std::string> g_callback_names;

static void attach_callback(const std::string &name);

// SharedString -> std::string
static std::string to_std(const slint::SharedString &s) {
  return std::string(s.data());
}

void slint_bridge_init(void) {
  if (!g_compiler)
    g_compiler = std::make_unique<ComponentCompiler>();
}

int slint_bridge_load(const char *slint_path) {
  slint_bridge_init();
  auto defOpt =
      g_compiler->build_from_path(slint_path ? slint_path : "ui/app.slint");
  if (!defOpt) {
    for (auto &d : g_compiler->diagnostics()) {
      fprintf(stderr, "Slint error: %s\n", to_std(d.message).c_str());
    }
    return -1;
  }
  auto def = *defOpt;
  g_instance = def.create();
  for (auto &name : g_callback_names)
    attach_callback(name);
  return 0;
}

void slint_bridge_run(void) {
  if (!g_instance.has_value()) {
    fprintf(stderr, "slint_bridge_run: pehle UI load karo\n");
    return;
  }
  (*g_instance)->show();
  slint::run_event_loop();
}

void slint_bridge_set_string(const char *name, const char *value) {
  if (!g_instance.has_value() || !name)
    return;
  (*g_instance)
      ->set_property(name, Value(slint::SharedString(value ? value : "")));
}
void slint_bridge_set_int(const char *name, int value) {
  if (!g_instance.has_value() || !name)
    return;
  (*g_instance)->set_property(name, Value(value));
}
void slint_bridge_set_bool(const char *name, int value) {
  if (!g_instance.has_value() || !name)
    return;
  (*g_instance)->set_property(name, Value(value != 0));
}
void slint_bridge_set_float(const char *name, double value) {
  if (!g_instance.has_value() || !name)
    return;
  (*g_instance)->set_property(name, Value(value));
}

// "#RRGGBB" or "#AARRGGBB" -> slint::Color
static slint::Color parse_hex_color(const char *hex) {
  std::string h = hex ? hex : "#000000";
  if (!h.empty() && h[0] == '#')
    h.erase(0, 1);
  unsigned long v = std::strtoul(h.c_str(), nullptr, 16);
  if (h.size() == 8) {
    return slint::Color::from_argb_uint8((v >> 24) & 0xFF, (v >> 16) & 0xFF,
                                         (v >> 8) & 0xFF, v & 0xFF);
  }
  return slint::Color::from_rgb_uint8((v >> 16) & 0xFF, (v >> 8) & 0xFF,
                                      v & 0xFF);
}

void slint_bridge_set_color(const char *name, const char *hex) {
  if (!g_instance.has_value() || !name)
    return;
  (*g_instance)->set_property(name, Value(slint::Brush(parse_hex_color(hex))));
}

void slint_bridge_set_image(const char *name, const char *path) {
  if (!g_instance.has_value() || !name || !path)
    return;
  (*g_instance)
      ->set_property(
          name, Value(slint::Image::load_from_path(slint::SharedString(path))));
}

void slint_bridge_set_image_bytes(const char *name, const unsigned char *data,
                                  int len) {
  if (!g_instance.has_value() || !name || !data || len <= 0)
    return;
  auto tmp = std::filesystem::temp_directory_path() / "slint_bridge_image.bin";
  {
    std::ofstream f(tmp, std::ios::binary);
    f.write(reinterpret_cast<const char *>(data), len);
  }
  slint_bridge_set_image(name, tmp.string().c_str());
}

char *slint_bridge_get_string(const char *name) {
  if (!g_instance.has_value() || !name)
    return strdup("");
  auto opt = (*g_instance)->get_property(name);
  if (!opt)
    return strdup("");
  auto s = opt->to_string();
  if (!s)
    return strdup("");
  std::string out = to_std(*s);
  return strdup(out.c_str());
}
int slint_bridge_get_int(const char *name) {
  if (!g_instance.has_value() || !name)
    return 0;
  auto opt = (*g_instance)->get_property(name);
  if (!opt)
    return 0;
  auto d = opt->to_number();
  return d ? (int)*d : 0;
}
int slint_bridge_get_bool(const char *name) {
  if (!g_instance.has_value() || !name)
    return 0;
  auto opt = (*g_instance)->get_property(name);
  if (!opt)
    return 0;
  auto b = opt->to_bool();
  return (b && *b) ? 1 : 0;
}

static void attach_callback(const std::string &name) {
  if (!g_instance.has_value())
    return;
  (*g_instance)->set_callback(name, [name](auto args) -> Value {
    std::string input;
    if (!args.empty()) {
      if (auto s = args[0].to_string())
        input = to_std(*s);
    }
    const char *out = GoCallbackHandler(name.c_str(), input.c_str());
    std::string res = out ? out : "";
    std::free((void *)out); // Go ne C.CString se banaya tha
    return Value(slint::SharedString(res));
  });
}

void slint_bridge_on_callback(const char *name) {
  if (!name)
    return;
  std::string n = name;
  for (auto &existing : g_callback_names) {
    if (existing == n) {
      attach_callback(n);
      return;
    }
  }
  g_callback_names.push_back(n);
  attach_callback(n);
}

// ---- Arrays ----
void slint_bridge_set_string_array(const char *name, const char **items,
                                   int n) {
  if (!g_instance || !name || !items || n < 0)
    return;
  slint::SharedVector<Value> vec;
  for (int i = 0; i < n; i++) {
    vec.push_back(Value(slint::SharedString(items[i] ? items[i] : "")));
  }
  (*g_instance)->set_property(name, Value(vec));
}

void slint_bridge_set_int_array(const char *name, const int *items, int n) {
  if (!g_instance || !name || !items || n < 0)
    return;
  slint::SharedVector<Value> vec;
  for (int i = 0; i < n; i++)
    vec.push_back(Value(items[i]));
  (*g_instance)->set_property(name, Value(vec));
}

void slint_bridge_set_float_array(const char *name, const double *items,
                                  int n) {
  if (!g_instance || !name || !items || n < 0)
    return;
  slint::SharedVector<Value> vec;
  for (int i = 0; i < n; i++)
    vec.push_back(Value(items[i]));
  (*g_instance)->set_property(name, Value(vec));
}

// ---- Structs ----
using SlintStructT = slint::interpreter::Struct;

void *slint_bridge_struct_new(void) { return new SlintStructT(); }

void slint_bridge_struct_free(void *s) {
  delete static_cast<SlintStructT *>(s);
}

void slint_bridge_struct_set_string(void *s, const char *field,
                                    const char *value) {
  if (!s || !field)
    return;
  static_cast<SlintStructT *>(s)->set_field(
      field, Value(slint::SharedString(value ? value : "")));
}

void slint_bridge_struct_set_int(void *s, const char *field, int value) {
  if (!s || !field)
    return;
  static_cast<SlintStructT *>(s)->set_field(field, Value(value));
}

void slint_bridge_struct_set_float(void *s, const char *field, double value) {
  if (!s || !field)
    return;
  static_cast<SlintStructT *>(s)->set_field(field, Value(value));
}

void slint_bridge_struct_set_bool(void *s, const char *field, int value) {
  if (!s || !field)
    return;
  static_cast<SlintStructT *>(s)->set_field(field, Value(value != 0));
}

void slint_bridge_struct_set_color(void *s, const char *field,
                                   const char *hex) {
  if (!s || !field)
    return;
  static_cast<SlintStructT *>(s)->set_field(
      field, Value(slint::Brush(parse_hex_color(hex))));
}

void slint_bridge_set_struct(const char *name, void *s) {
  if (!g_instance || !name || !s)
    return;
  (*g_instance)->set_property(name, Value(*static_cast<SlintStructT *>(s)));
}

void slint_bridge_set_struct_array(const char *name, void **structs, int n) {
  if (!g_instance || !name || !structs || n < 0)
    return;
  slint::SharedVector<Value> vec;
  for (int i = 0; i < n; i++) {
    auto *s = static_cast<SlintStructT *>(structs[i]);
    if (s)
      vec.push_back(Value(
          *s)); // struct copy hota hai, Go ka struct baad me free ho sakta hai
  }
  (*g_instance)->set_property(name, Value(vec));
}
