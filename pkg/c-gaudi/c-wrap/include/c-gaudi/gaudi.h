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


/* StatusCode */
struct CGaudi_StatusCode {
  unsigned long   code;      ///< The status code
};

typedef void* CGaudi_IInterface;
typedef void* CGaudi_INamedInterface;
typedef void* CGaudi_IAlgorithm;
typedef void* CGaudi_IService;
typedef void* CGaudi_IAlgTool;
typedef void* CGaudi_IApplicationMgr;

/* InterfaceID */

  struct CGaudi_InterfaceID {
    unsigned long id;
    unsigned long major_ver;
    unsigned long minor_ver;                                                                                                                                                                                                                                                 
  };

CGAUDI_API
int
CGaudi_InterfaceID_versionMatch(CGaudi_InterfaceID self, CGaudi_InterfaceID other);

CGAUDI_API
int
CGaudi_InterfaceID_fullMatch(CGaudi_InterfaceID self, CGaudi_InterfaceID other);

/* IInterface */

CGAUDI_API
CGaudi_InterfaceID
CGaudi_IInterface_InterfaceID(CGaudi_IInterface self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IInterface_queryInterface(CGaudi_IInterface self, CGaudi_InterfaceID iid, void **p);

CGAUDI_API
unsigned long
CGaudi_IInterface_addRef(CGaudi_IInterface self);

CGAUDI_API
unsigned long
CGaudi_IInterface_release(CGaudi_IInterface self);

CGAUDI_API
unsigned long
CGaudi_IInterface_refCount(CGaudi_IInterface self);

/* INamedInterface */
CGAUDI_API
const char*
CGaudi_INamedInterface_name(CGaudi_INamedInterface self);

/* IAlgorithm */
CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgorithm_execute(CGaudi_IAlgorithm self);

CGAUDI_API
int 
CGaudi_IAlgorithm_isInitialized(CGaudi_IAlgorithm self);

CGAUDI_API
int
CGaudi_IAlgorithm_isFinalized(CGaudi_IAlgorithm self);

CGAUDI_API
int
CGaudi_IAlgorithm_isExecuted(CGaudi_IAlgorithm self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgorithm_configure(CGaudi_IAlgorithm self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgorithm_initialize(CGaudi_IAlgorithm self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgorithm_start(CGaudi_IAlgorithm self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgorithm_stop(CGaudi_IAlgorithm self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgorithm_finalized(CGaudi_IAlgorithm self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgorithm_terminate(CGaudi_IAlgorithm self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgorithm_reinitialize(CGaudi_IAlgorithm self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgorithm_restart(CGaudi_IAlgorithm self);


  /** Get the current state.
   */
  /* virtual Gaudi::StateMachine::State FSMState() const = 0; */

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgorithm_sysInitialize(CGaudi_IAlgorithm self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgorithm_sysReinitialize(CGaudi_IAlgorithm self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgorithm_sysRestart(CGaudi_IAlgorithm self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgorithm_sysExecute(CGaudi_IAlgorithm self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgorithm_sysStop(CGaudi_IAlgorithm self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgorithm_sysFinalize(CGaudi_IAlgorithm self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgorithm_sysBeginRun(CGaudi_IAlgorithm self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgorithm_sysEndRun(CGaudi_IAlgorithm self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgorithm_beginRun(CGaudi_IAlgorithm self);

CGAUDI_API
CGaudi_StatusCode
CGaudi_IAlgorithm_endRun(CGaudi_IAlgorithm self);

CGAUDI_API
void
CGaudi_IAlgorithm_resetExecuted(CGaudi_IAlgorithm self);

CGAUDI_API
void
CGaudi_IAlgorithm_setExecuted(CGaudi_IAlgorithm self, int state);

CGAUDI_API
int
CGaudi_IAlgorithm_isEnabled(CGaudi_IAlgorithm self);

CGAUDI_API
int
CGaudi_IAlgorithm_filterPassed(CGaudi_IAlgorithm self);

CGAUDI_API
void
CGaudi_IAlgorithm_setFilterPassed(CGaudi_IAlgorithm self, int state);

#ifdef __cplusplus
} /* !extern "C" */
#endif

#endif /*!CGAUDI_GAUDI_H */
