#include "c-gaudi/gaudi.h"

#include "GaudiKernel/IInterface.h"

CGaudi_InterfaceID
CGaudi_IInterface_InterfaceID(CGaudi_IInterface self)
{
  return (CGaudi_InterfaceID)((IInterface*)self)->interfaceID();
}
