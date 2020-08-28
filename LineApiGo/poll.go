// Copyright (c) 2020 @Ch31212y
// Version 1.1 beta
// LastUpdate 2020/08/28

package lineapigo

import ser "github.com/ch31212y/lineapigo/talkservice"

// FetchOps fetch operations
func (cl *LineClient) FetchOps() ([]*ser.Operation, error) {
	return cl.poll.FetchOps(cl.ctx, cl.revision, 100, cl.globalRev, cl.individualRev)
}

// FetchOperations fetch operations
func (cl *LineClient) FetchOperations() ([]*ser.Operation, error) {
	return cl.poll.FetchOperations(cl.ctx, cl.revision, 100)
}

// SetRevision set revision
func (cl *LineClient) SetRevision(rev int64) {
	if rev > cl.revision {
		cl.revision = rev
	}
}
