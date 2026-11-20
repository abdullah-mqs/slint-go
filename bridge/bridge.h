#pragma once
// Generic API, not Hard coded
#ifdef __cplusplus
extern "C" {
#endif
#include <stddef.h>

void slint_bridge_init(void);
// 0 = success, -1 = UI load/compile fail
int slint_bridge_load(const char *slint_path);
void slint_bridge_run(void);

void slint_bridge_set_string(const char *name, const char *value);
void slint_bridge_set_int(const char *name, int value);
void slint_bridge_set_bool(const char *name, int value);
void slint_bridge_set_float(const char *name, double value);
void slint_bridge_set_color(const char *name,
                            const char *hex); // "#FF5733" ya "#FF5733AA"
void slint_bridge_set_image(const char *name, const char *path);
void slint_bridge_set_image_bytes(const char *name, const unsigned char *data,
                                  int len);

char *slint_bridge_get_string(const char *name); // caller free() kare
int slint_bridge_get_int(const char *name);
int slint_bridge_get_bool(const char *name);

// Go callback
const char *GoCallbackHandler(const char *name, const char *input);
void slint_bridge_on_callback(const char *name);

// Arrays
void slint_bridge_set_string_array(const char *name, const char **items, int n);
void slint_bridge_set_int_array(const char *name, const int *items, int n);
void slint_bridge_set_float_array(const char *name, const double *items, int n);

// Structs (opaque handle)
void *slint_bridge_struct_new(void);
void slint_bridge_struct_free(void *s);
void slint_bridge_struct_set_string(void *s, const char *field,
                                    const char *value);
void slint_bridge_struct_set_int(void *s, const char *field, int value);
void slint_bridge_struct_set_float(void *s, const char *field, double value);
void slint_bridge_struct_set_bool(void *s, const char *field, int value);
void slint_bridge_struct_set_color(void *s, const char *field, const char *hex);
void slint_bridge_set_struct(const char *name, void *s);
void slint_bridge_set_struct_array(const char *name, void **structs, int n);
#ifdef __cplusplus
}
#endif
