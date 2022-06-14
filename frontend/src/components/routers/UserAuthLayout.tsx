import React from "react";
import { Navigate, Outlet } from "react-router-dom";
import { useAuth } from "../contexts/AuthContext";

export default function UserAuthLayout() {
  const { state } = useAuth();

  return (state.isAuthorized) ? <Outlet /> : <Navigate to="/auth/sign_in" />;
}
