#include "c-gaudi/gaudi.h"

#include "GaudiKernel/IInterface.h"
#include "GaudiKernel/INamedInterface.h"
#include "GaudiKernel/IAlgorithm.h"

/* IInterface */

CGaudi_InterfaceID
CGaudi_IInterface_InterfaceID(CGaudi_IInterface self)
{
  const InterfaceID& iid = ((IInterface*)self)->interfaceID();
  return *(CGaudi_InterfaceID*)(&iid);
}

int
CGaudi_InterfaceID_versionMatch(CGaudi_InterfaceID self, CGaudi_InterfaceID other)
{
  bool rc = ((InterfaceID*)&self)->versionMatch(*(InterfaceID*)&other);
  return rc ? 1 : 0;
}

int
CGaudi_InterfaceID_fullMatch(CGaudi_InterfaceID self, CGaudi_InterfaceID other)
{
  bool rc = ((InterfaceID*)&self)->fullMatch(*(InterfaceID*)&other);
  return rc ? 1 : 0;
}

/* IInterface */

CGaudi_StatusCode
CGaudi_IInterface_queryInterface(CGaudi_IInterface self, CGaudi_InterfaceID iid, void **p)
{
  StatusCode sc = ((IInterface*)self)->queryInterface(*(InterfaceID*)&iid, p);
  return *(CGaudi_StatusCode*)(&sc);
}

unsigned long
CGaudi_IInterface_addRef(CGaudi_IInterface self)
{
  return ((IInterface*)self)->addRef();
}

unsigned long
CGaudi_IInterface_release(CGaudi_IInterface self)
{
  return ((IInterface*)self)->release();
}

unsigned long
CGaudi_IInterface_refCount(CGaudi_IInterface self)
{
  return ((IInterface*)self)->refCount();
}

/* INamedInterface */

const char*
CGaudi_INamedInterface_name(CGaudi_INamedInterface self)
{
  return ((INamedInterface*)self)->name().c_str();
}


/* IAlgorithm */

CGaudi_StatusCode
CGaudi_IAlgorithm_execute(CGaudi_IAlgorithm self)
{
  StatusCode sc = ((IAlgorithm*)self)->execute();
  return *(CGaudi_StatusCode*)(&sc);
}


int 
CGaudi_IAlgorithm_isInitialized(CGaudi_IAlgorithm self)
{
  return ((IAlgorithm*)self)->isInitialized() ? 1 : 0;
}


int
CGaudi_IAlgorithm_isFinalized(CGaudi_IAlgorithm self)
{
  return ((IAlgorithm*)self)->isFinalized() ? 1 : 0;
}


int
CGaudi_IAlgorithm_isExecuted(CGaudi_IAlgorithm self)
{
  return ((IAlgorithm*)self)->isExecuted() ? 1 : 0;
}


CGaudi_StatusCode
CGaudi_IAlgorithm_configure(CGaudi_IAlgorithm self)
{
  StatusCode sc = ((IAlgorithm*)self)->configure();
  return *(CGaudi_StatusCode*)(&sc);
}

CGaudi_StatusCode
CGaudi_IAlgorithm_initialize(CGaudi_IAlgorithm self)
{
  StatusCode sc = ((IAlgorithm*)self)->initialize();
  return *(CGaudi_StatusCode*)(&sc);
}

CGaudi_StatusCode
CGaudi_IAlgorithm_start(CGaudi_IAlgorithm self)
{
  StatusCode sc = ((IAlgorithm*)self)->start();
  return *(CGaudi_StatusCode*)(&sc);
}

CGaudi_StatusCode
CGaudi_IAlgorithm_stop(CGaudi_IAlgorithm self)
{
  StatusCode sc = ((IAlgorithm*)self)->stop();
  return *(CGaudi_StatusCode*)(&sc);
}

CGaudi_StatusCode
CGaudi_IAlgorithm_finalize(CGaudi_IAlgorithm self)
{
  StatusCode sc = ((IAlgorithm*)self)->finalize();
  return *(CGaudi_StatusCode*)(&sc);
}

CGaudi_StatusCode
CGaudi_IAlgorithm_terminate(CGaudi_IAlgorithm self)
{
  StatusCode sc = ((IAlgorithm*)self)->terminate();
  return *(CGaudi_StatusCode*)(&sc);
}

CGaudi_StatusCode
CGaudi_IAlgorithm_reinitialize(CGaudi_IAlgorithm self)
{
  StatusCode sc = ((IAlgorithm*)self)->reinitialize();
  return *(CGaudi_StatusCode*)(&sc);
}

CGaudi_StatusCode
CGaudi_IAlgorithm_restart(CGaudi_IAlgorithm self)
{
  StatusCode sc = ((IAlgorithm*)self)->restart();
  return *(CGaudi_StatusCode*)(&sc);
}

  /** Get the current state.
   */
  /* virtual Gaudi::StateMachine::State FSMState() const = 0; */


CGaudi_StatusCode
CGaudi_IAlgorithm_sysInitialize(CGaudi_IAlgorithm self)
{
  StatusCode sc = ((IAlgorithm*)self)->sysInitialize();
  return *(CGaudi_StatusCode*)(&sc);
}

CGaudi_StatusCode
CGaudi_IAlgorithm_sysReinitialize(CGaudi_IAlgorithm self)
{
  StatusCode sc = ((IAlgorithm*)self)->sysReinitialize();
  return *(CGaudi_StatusCode*)(&sc);
}

CGaudi_StatusCode
CGaudi_IAlgorithm_sysRestart(CGaudi_IAlgorithm self)
{
  StatusCode sc = ((IAlgorithm*)self)->sysRestart();
  return *(CGaudi_StatusCode*)(&sc);
}

CGaudi_StatusCode
CGaudi_IAlgorithm_sysExecute(CGaudi_IAlgorithm self)
{
  StatusCode sc = ((IAlgorithm*)self)->sysExecute();
  return *(CGaudi_StatusCode*)(&sc);
}

CGaudi_StatusCode
CGaudi_IAlgorithm_sysStop(CGaudi_IAlgorithm self)
{
  StatusCode sc = ((IAlgorithm*)self)->sysStop();
  return *(CGaudi_StatusCode*)(&sc);
}

CGaudi_StatusCode
CGaudi_IAlgorithm_sysFinalize(CGaudi_IAlgorithm self)
{
  StatusCode sc = ((IAlgorithm*)self)->sysFinalize();
  return *(CGaudi_StatusCode*)(&sc);
}

CGaudi_StatusCode
CGaudi_IAlgorithm_sysBeginRun(CGaudi_IAlgorithm self)
{
  StatusCode sc = ((IAlgorithm*)self)->sysBeginRun();
  return *(CGaudi_StatusCode*)(&sc);
}

CGaudi_StatusCode
CGaudi_IAlgorithm_sysEndRun(CGaudi_IAlgorithm self)
{
  StatusCode sc = ((IAlgorithm*)self)->sysEndRun();
  return *(CGaudi_StatusCode*)(&sc);
}

CGaudi_StatusCode
CGaudi_IAlgorithm_beginRun(CGaudi_IAlgorithm self)
{
  StatusCode sc = ((IAlgorithm*)self)->beginRun();
  return *(CGaudi_StatusCode*)(&sc);
}

CGaudi_StatusCode
CGaudi_IAlgorithm_endRun(CGaudi_IAlgorithm self)
{
  StatusCode sc = ((IAlgorithm*)self)->endRun();
  return *(CGaudi_StatusCode*)(&sc);
}

void
CGaudi_IAlgorithm_resetExecuted(CGaudi_IAlgorithm self)
{
  ((IAlgorithm*)self)->resetExecuted();
}

void
CGaudi_IAlgorithm_setExecuted(CGaudi_IAlgorithm self, int state)
{
  ((IAlgorithm*)self)->setExecuted(state ? true : false);
}

int
CGaudi_IAlgorithm_isEnabled(CGaudi_IAlgorithm self)
{
  return ((IAlgorithm*)self)->isEnabled() ? 1 : 0;
}


int
CGaudi_IAlgorithm_filterPassed(CGaudi_IAlgorithm self)
{
  return ((IAlgorithm*)self)->filterPassed() ? 1 : 0;
}


void
CGaudi_IAlgorithm_setFilterPassed(CGaudi_IAlgorithm self, int state)
{
  ((IAlgorithm*)self)->setFilterPassed(state ? true : false);
}





