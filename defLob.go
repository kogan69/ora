// Copyright 2014 Rana Ian. All rights reserved.
// Use of this source code is governed by The MIT License
// found in the accompanying LICENSE file.

package ora

/*
#include <oci.h>
*/
import "C"
import (
	"github.com/golang/glog"
	"unsafe"
)

type defLob struct {
	rset          *Rset
	ocidef        *C.OCIDefine
	ociLobLocator *C.OCILobLocator
	null          C.sb2
	gct           GoColumnType
	sqlt          C.ub2
	charsetForm   C.ub1
}

func (def *defLob) define(position int, charsetForm C.ub1, sqlt C.ub2, gct GoColumnType, rset *Rset) error {
	glog.Infoln("position: ", position)
	def.rset = rset
	def.gct = gct
	def.sqlt = sqlt
	def.charsetForm = charsetForm
	r := C.OCIDefineByPos2(
		def.rset.ocistmt,                        //OCIStmt     *stmtp,
		&def.ocidef,                             //OCIDefine   **defnpp,
		def.rset.stmt.ses.srv.env.ocierr,        //OCIError    *errhp,
		C.ub4(position),                         //ub4         position,
		unsafe.Pointer(&def.ociLobLocator),      //void        *valuep,
		C.sb8(unsafe.Sizeof(def.ociLobLocator)), //sb8         value_sz,
		sqlt, //ub2         dty,
		unsafe.Pointer(&def.null), //void        *indp,
		nil,           //ub2         *rlenp,
		nil,           //ub2         *rcodep,
		C.OCI_DEFAULT) //ub4         mode );
	if r == C.OCI_ERROR {
		return def.rset.stmt.ses.srv.env.ociError()
	}
	return nil
}
func (def *defLob) Bytes() (value []byte, err error) {
	var lobLength C.oraub8
	// Open the lob to obtain length; round-trip to database
	r := C.OCILobOpen(
		def.rset.stmt.ses.srv.ocisvcctx,  //OCISvcCtx          *svchp,
		def.rset.stmt.ses.srv.env.ocierr, //OCIError           *errhp,
		def.ociLobLocator,                //OCILobLocator      *locp,
		C.OCI_LOB_READONLY)               //ub1              mode );
	if r == C.OCI_ERROR {
		return nil, def.rset.stmt.ses.srv.env.ociError()
	}
	// get the length of the lob
	r = C.OCILobGetLength2(
		def.rset.stmt.ses.srv.ocisvcctx,  //OCISvcCtx          *svchp,
		def.rset.stmt.ses.srv.env.ocierr, //OCIError           *errhp,
		def.ociLobLocator,                //OCILobLocator      *locp,
		&lobLength)                       //oraub8 *lenp)
	if r == C.OCI_ERROR {
		return nil, def.rset.stmt.ses.srv.env.ociError()
	}

	if lobLength > 0 {
		// Allocate []byte the length of the lob
		value = make([]byte, int(lobLength))
		// buffer is size of ora.LobBufferSize
		var buffer [1 << 24]byte
		var writeIndex int
		var byte_amtp C.oraub8 = lobLength
		var piece C.ub1 = C.OCI_FIRST_PIECE
		var loading bool = true
		for loading {
			r = C.OCILobRead2(
				def.rset.stmt.ses.srv.ocisvcctx,  //OCISvcCtx          *svchp,
				def.rset.stmt.ses.srv.env.ocierr, //OCIError           *errhp,
				def.ociLobLocator,                //OCILobLocator      *locp,
				&byte_amtp,                       //oraub8             *byte_amtp,
				nil,                              //oraub8             *char_amtp,
				C.oraub8(1),                      //oraub8             offset, offset is 1-based
				unsafe.Pointer(&buffer[0]),       //void               *bufp,
				C.oraub8(len(buffer)),            //oraub8             bufl,
				piece,           //ub1                piece,
				nil,             //void               *ctxp,
				nil,             //OCICallbackLobRead2 (cbfp)
				C.ub2(0),        //ub2                csid,
				def.charsetForm) //ub1                csfrm );

			if r == C.OCI_ERROR {
				return nil, def.rset.stmt.ses.srv.env.ociError()
			} else {
				// Write buffer to return slice
				// byte_amtp represents the amount copied into buffer by oci
				for n := 0; n < int(byte_amtp); n++ {
					value[writeIndex] = buffer[n]
					writeIndex++
				}
				// Determine action for next cycle
				if r == C.OCI_NEED_DATA {
					piece = C.OCI_NEXT_PIECE
				} else if r == C.OCI_SUCCESS {
					loading = false
				}
			}
		}
	}

	r = C.OCILobClose(
		def.rset.stmt.ses.srv.ocisvcctx,  //OCISvcCtx          *svchp,
		def.rset.stmt.ses.srv.env.ocierr, //OCIError           *errhp,
		def.ociLobLocator)                //OCILobLocator      *locp,
	if r == C.OCI_ERROR {
		return nil, def.rset.stmt.ses.srv.env.ociError()
	}

	return value, nil

}
func (def *defLob) String() (value string, err error) {
	var bytes []byte
	bytes, err = def.Bytes()
	value = string(bytes)
	return value, err
}
func (def *defLob) value() (value interface{}, err error) {
	if def.sqlt == C.SQLT_BLOB {
		if def.gct == Bin {
			if def.null > -1 {
				value, err = def.Bytes()
			}
		} else {
			binValue := Binary{IsNull: def.null < 0}
			if !binValue.IsNull {
				binValue.Value, err = def.Bytes()
			}
			value = binValue
		}
	} else {
		if def.gct == S {
			if def.null > -1 {
				value, err = def.String()
			}
		} else {
			oraString := String{IsNull: def.null < 0}
			if !oraString.IsNull {
				oraString.Value, err = def.String()
			}
			value = oraString
		}
	}
	return value, err
}
func (def *defLob) alloc() error {
	// Allocate lob locator handle
	// OCI_DTYPE_LOB is for a BLOB or CLOB
	r := C.OCIDescriptorAlloc(
		unsafe.Pointer(def.rset.stmt.ses.srv.env.ocienv),      //CONST dvoid   *parenth,
		(*unsafe.Pointer)(unsafe.Pointer(&def.ociLobLocator)), //dvoid         **descpp,
		C.OCI_DTYPE_LOB,                                       //ub4           type,
		0,                                                     //size_t        xtramem_sz,
		nil)                                                   //dvoid         **usrmempp);
	if r == C.OCI_ERROR {
		return def.rset.stmt.ses.srv.env.ociError()
	} else if r == C.OCI_INVALID_HANDLE {
		return errNew("unable to allocate oci lob handle during define")
	}
	return nil
}

func (def *defLob) free() {
	defer func() {
		recover()
	}()
	C.OCIDescriptorFree(
		unsafe.Pointer(def.ociLobLocator), //void     *descp,
		C.OCI_DTYPE_LOB)                   //ub4      type );
}

func (def *defLob) close() (err error) {
	defer func() {
		if value := recover(); value != nil {
			err = errRecover(value)
		}
	}()

	glog.Infoln("close")
	def.free()
	rset := def.rset
	def.rset = nil
	def.ocidef = nil
	def.ociLobLocator = nil
	rset.putDef(defIdxLob, def)
	return nil
}
