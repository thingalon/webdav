package webdav

import (
	"fmt"
	"strings"
	"time"
)

type Lock struct {
	uri      string
	creator  string
	owner    string
	depth    int
	timeout  time.Duration
	typ      string
	scope    string
	token    string
	Modified time.Time
}

func NewLock(uri, creator, owner string) *Lock {
	return &Lock{
		uri,
		creator,
		owner,
		0,
		0,
		"write",
		"exclusive",
		generateToken(),
		time.Now(),
	}
}

func (lock *Lock) Refresh() {
	lock.Modified = time.Now()
}

func (lock *Lock) IsValid() bool {
	return lock.timeout > time.Now().Sub(lock.Modified)
}

func (lock *Lock) GetTimeoutString() string {
	return fmt.Sprintf("Second-%d", lock.timeout/time.Second)
}

func (lock *Lock) setTimeout(timeout time.Duration) {
	lock.timeout = timeout
	lock.Modified = time.Now()
}

func (lock *Lock) asXML(namespace string, discover bool) string {
	//owner_str = lock.owner
	//owner_str = "".join([node.toxml() for node in self.owner[0].childNodes])

	base := fmt.Sprintf(`<%[1]s:activelock>
             <%[1]s:locktype><%[1]s:%[2]s/></%[1]s:locktype>
             <%[1]s:lockscope><%[1]s:%[3]s/></%[1]s:lockscope>
             <%[1]s:depth>%[4]d</%[1]s:depth>
             <%[1]s:owner>%[5]s</%[1]s:owner>
             <%[1]s:timeout>%[6]s</%[1]s:timeout>
             <%[1]s:locktoken>
             <%[1]s:href>opaquelocktoken:%[7]s</%[1]s:href>
             </%[1]s:locktoken>
             </%[1]s:activelock>
             `, strings.Trim(namespace, ":"),
		lock.typ,
		lock.scope,
		lock.depth,
		lock.owner,
		lock.GetTimeoutString(),
		lock.token,
	)

	if discover {
		return base
	}

	return fmt.Sprintf(`<?xml version="1.0" encoding="utf-8" ?>
<d:prop xmlns:d="DAV:">
 <d:lockdiscovery>
  %s
 </d:lockdiscovery>
</d:prop>`, base)
}
