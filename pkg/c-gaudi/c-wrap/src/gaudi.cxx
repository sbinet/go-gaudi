#include "c-gaudi/gaudi.h"

#include "GaudiKernel/IInterface.h"
#include "GaudiKernel/INamedInterface.h"

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


