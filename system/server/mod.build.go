package server

import (
	"github/anurag/altar-be/modules/auth"
)

func (s *Server) BuildModules() {
	auth.BuildAuthModule(&s.ROUTER.RouterGroup, s.PG_CLIENT, s.RDS_CLIENT, s.CONFIG, s.JWT_MANAGER, s.GOOGLE_PROVIDER)
}
