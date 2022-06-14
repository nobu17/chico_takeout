import React from "react";
import { Outlet, Navigate } from "react-router-dom";
import { useAuth } from "../contexts/AuthContext";

export default function UserLoginLayout() {
  const { state } = useAuth();

  return (state.isAuthorized) ? <Navigate to="/my_page" /> : <Outlet />;
}
