package overmind

import (
	"time"
)

const (
	MAX_BODY_SIZE  = 1048576
	NICE_MAGIC_NUM = 0xFEA94831
)

const (
	DEFAULT_SEND_BUFFER_SIZE    = 15
	DEFAULT_RECEIVE_BUFFER_SIZE = 15
)

const (
	DEFAULT_FIRST_SEQNUM       = 0
	DEFAULT_HEARTBEAT_TIME_OUT = 180
)

const (
	DEFAULT_WRITE_TIME_OUT = time.Second * 60
	DEFAULT_READ_TIME_OUT  = time.Second * (DEFAULT_HEARTBEAT_TIME_OUT + 60)
	DEFAULT_AUTH_TIME_OUT  = time.Second * 10
	MAX_ERROR_CNT          = 5
)

const (
	DEFAULT_KEEPALIVE_PERIOD = time.Second * DEFAULT_HEARTBEAT_TIME_OUT * 2

	// to check dead-conn
	PERIOD_NO_ACK_CNT      = 4
	PERIOD_NO_ACK_TIME_OUT = time.Second * 120
	MAX_NO_ACK_CNT         = 12
)

const (
	// for retransmission
	TIMERQUEUE_CLUSTER_NODE_NUM = 500
	TIMERQUEUE_DELAYINTERVAL    = time.Second * 60
)

const (
	HANDSHAKE_TYPE = 254
	HEARTBEAT_TYPE = 255
)

const (
	REDIS_KEY_TCP_CONNECT             = "overmind_tcp_connect:%s"
	REDIS_KEY_USER_TCP_CONNECT_ADDR   = "overmind_user_tcp_connect_addr:%s"
	REDIS_KEY_USER_PROXY_CONNECT_ADDR = "overmind_user_proxy_connect_addr:%s"
	REDIS_KEY_USER_TOKEN_CONNECT      = "overmind_user_token_connect:%s"
	REDIS_KEY_USER_CONNECT_STATUS     = "overmind_user_connect_status:%s"

	REDIS_KEY_USER_WEBSOCKET_CONNECT = "overmind_user_websocket_connect:%s"

	ACTIVE_REDIS_KEY_USER_STATUS = "user_activity_time_%s"
)

const (
	DEFAULT_UPDATE_INTERVAL = time.Minute * 5
)
