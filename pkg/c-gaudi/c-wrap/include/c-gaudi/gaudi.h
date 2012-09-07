/* c-gaudi */
#ifndef CGAUDI_GAUDI_H
#define CGAUDI_GAUDI_H 1

#include <stdint.h>
#include <stddef.h>
#include <stdbool.h>

#if __GNUC__ >= 4
#  define CGAUDI_HASCLASSVISIBILITY
#endif

#if defined(CGAUDI_HASCLASSVISIBILITY)
#  define CGAUDI_IMPORT __attribute__((visibility("default")))
#  define CGAUDI_EXPORT __attribute__((visibility("default")))
#  define CGAUDI_LOCAL  __attribute__((visibility("hidden")))
#else
#  define CGAUDI_IMPORT
#  define CGAUDI_EXPORT
#  define CGAUDI_LOCAL
#endif

#define CGAUDI_API CGAUDI_EXPORT

#ifdef __cplusplus
extern "C" {
#endif

typedef void* CGaudi_IAlgorithm;
typedef void* CGaudi_IService;
typedef void* CGaudi_IAlgTool;
typedef void* CGaudi_IApplicationMgr;
typedef void* CGaudi_IInterface;
typedef void* CGaudi_INamedInterface;

struct CGaudi_InterfaceID {
  unsigned long id;
  unsigned long major_ver;
  unsigned long minor_ver;
};

CGAUDI_API
CGaudi_InterfaceID
CGaudi_IInterface_InterfaceID(CGaudi_IInterface self);

#ifdef __cplusplus
} /* !extern "C" */
#endif

#endif /*!CGAUDI_GAUDI_H */
