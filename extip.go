package extip

import (
    "gortc.io/stun"
    "context"
    "errors"
)

var PublicServerList = []string{
    "stun.l.google.com:19302",
    "stun.ekiga.net:3478",
    "stun.ideasip.com:3478",
    "stun.schlund.de:3478",
    "stun.voiparound.com:3478",
    "stun.voipbuster.com:3478",
    "stun.voipstunt.com:3478",
}

// Query external IP address from single server
func QuerySingleServer(ctx context.Context, server string, ipv6 bool) (string, error) {
    family := "udp4"
    if ipv6 {
        family = "udp6"
    }
    c, err := stun.Dial(family, server)
    if err != nil {
        return "", err
    }
    defer c.Close()

    message, err := stun.Build(stun.TransactionID, stun.BindingRequest)
    if err != nil {
        return "", err
    }

    clientOut := make(chan stun.Event)
    clientErr := make(chan error)

    go func() {
        err := c.Do(message, func(res stun.Event) {
            clientOut <- res
        })
        if err != nil {
            clientErr <- err
        }
    }()

    select {
    case err := <-clientErr:
        return "", err
    case res := <-clientOut:
		if res.Error != nil {
            return "", res.Error
		}
		var xorAddr stun.XORMappedAddress
		if err := xorAddr.GetFrom(res.Message); err == nil {
            return xorAddr.IP.String(), nil
		} else {
            var mappedAddr stun.MappedAddress
            if err := mappedAddr.GetFrom(res.Message); err == nil {
                return mappedAddr.IP.String(), nil
            } else {
                return "", err
            }
        }
    case <-ctx.Done():
        return "", errors.New("context was cancelled")
    }
}

// Query multiple servers and determine result by quorum of successful responses.
// List of servers can be a nil slice. In this case public server list will be used.
func QueryMultipleServers(ctx context.Context, servers []string, quorum uint, ipv6 bool) (string, error) {
    if servers == nil {
        servers = PublicServerList
    }
    count := len(servers)
    if count == 0 {
        return "", errors.New("empty server list")
    }
    if uint(len(servers)) < quorum {
        return "", errors.New("quorum is higher than server list length")
    }

    clientRes := make(chan string)
    clientErr := make(chan error)
    for _, server := range servers {
        go func(srv string) {
            res, err := QuerySingleServer(ctx, srv, ipv6)
            if err != nil {
                clientErr <- err
            } else {
                clientRes <- res
            }
        }(server)
    }

    resultMap := make(map[string]uint)
    var errorList []error

    for i := 0; i < count; i++ {
        select {
        case err := <-clientErr:
            errorList = append(errorList, err)
        case res := <-clientRes:
            resultMap[res]++
            if resultMap[res] >= quorum {
                return res, nil
            }
        }
    }
    return "", errors.New("no conclusive result was retrieved")
}
