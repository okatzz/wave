package engine

import (
	"context"

	"github.com/immesys/wave/iapi"
)

func (e *Engine) getEntityFromHashLoc(ctx context.Context, hash iapi.HashSchemeInstance, loc iapi.LocationSchemeInstance) (*iapi.Entity, error) {
	panic("ni")
}

//Check for revocations
//If decoded, check for expiry and entRevoked
//This must be fairly fast as it gets called frequently on the same stuff
func (e *Engine) checkAttestationAndSave(ctx context.Context, d *iapi.Attestation) (okay bool, err error) {
	attesterh, attloc, err := d.Attester()
	if err != nil {
		return false, err
	}
	attester, err := getEntityFromHashLoc(ctx, attesterh, attloc)
	if err != nil {
		return false, err
	}
	srcokay, err := e.checkEntityAndSave(attester)
	if err != nil {
		return false, err
	}
	subjecth, subjloc := d.Subject()
	subject, err := getEntityFromHashLoc(ctx, subjecth, subjloc)
	if err != nil {
		return false, err
	}
	dstokay, err := e.checkEntityAndSave(subject)
	if err != nil {
		return false, err
	}
	expired, err := d.Expired()
	if err != nil {
		return false, err
	}
	//TODO
	/*
		revoked, err := e.IsRevoked(e.ctx, d.PlaintextHeader.RevocationHash)
		if err != nil {
			return false, err
		}
	*/

	if !srcokay || !dstokay {
		//This dot must move to EntRevoked
		return false, e.ws.MoveAttestationEntRevoked(e.ctx, d)
	}
	if revoked {
		return false, e.ws.MoveAttestationRevokedG(e.ctx, d)
	}
	if expired {
		return false, e.ws.MoveAttestationExpiredP(e.ctx, d)
	}
	return true, nil
}

func (e *Engine) checkPendingAttestationAndSave(d *iapi.Attestation) (okay bool, err error) {
	//Like checkDot but don't check (nonexistant) content
	subjecth, subjloc := d.Subject()
	subject, err := getEntityFromHashLoc(ctx, subjecth, subjloc)
	if err != nil {
		return false, err
	}
	dstokay, err := e.checkEntityAndSave(subject)
	if err != nil {
		return false, err
	}
	//TODO
	/*
		revoked, err := e.IsRevoked(e.ctx, d.PlaintextHeader.RevocationHash)
		if err != nil {
			return false, err
		}
	*/
	if !dstokay {
		//This dot must move to EntRevoked
		return false, e.ws.MoveAttestationEntRevokedP(e.ctx, d)
	}
	// if revoked {
	// 	return false, e.ws.MoveDotRevokedG(e.ctx, d)
	// }
	return true, nil
}

func (e *Engine) checkEntityAndSave(ent *iapi.Entity) (bool, error) {
	if ent.Expired() {
		return false, e.ws.MoveEntityExpiredG(e.ctx, ent)
	}
	//TODO
	/*
		revoked, err := e.IsRevoked(e.ctx, ent.RevocationHash)
		if err != nil {
			return false, err
		}
		if revoked {
			return false, e.ws.MoveEntityRevokedG(e.ctx, ent)
		}
	*/
	return true, nil
}
