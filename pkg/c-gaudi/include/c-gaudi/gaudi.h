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

void *CGaudi_IAlgorithm;
void *CGaudi_IService;
void *CGaudi_IAlgTool;
void *CGaudi_IApplicationMgr;
void *CGaudi_IInterface;
void *CGaudi_INamedInterface;

struct CGaudi_InterfaceID {
  unsigned long m_id;
  unsigned long m_major_ver;
  unsigned long m_minor_ver;
};

CGAUDI_API
CGaudi_InterfaceID
CGaudi_IInterface_InterfaceID(CGaudi_IInterface self);

#ifdef __cplusplus
} /* !extern "C" */
#endif

#endif /*!CGAUDI_GAUDI_H */
