#include "c-gaudi/gaudi.h"

#include "GaudiKernel/IInterface.h"

CGaudi_InterfaceID
CGaudi_IInterface_InterfaceID(CGaudi_IInterface self)
{
  InterfaceID id = ((IInterface*)self)->interfaceID();
  CGaudi_InterfaceID cid = {
    id.id(),
    id.majorVersion(),
    id.minorVersion()
  };
  
  return cid;
}
