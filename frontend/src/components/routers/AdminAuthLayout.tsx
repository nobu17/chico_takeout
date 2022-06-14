import React from "react";
import { Navigate, Outlet } from "react-router-dom";
import { useAuth } from "../contexts/AuthContext";

export default function AdminAuthLayout() {
  const { state } = useAuth();
  if (state.isAuthorized && state.isAdmin) {
    return <Outlet />;
  }
  if (state.isAuthorized) {
    return <Navigate to="/" />;
  }
  return <Navigate to="/admin/sign_in" />;
}
