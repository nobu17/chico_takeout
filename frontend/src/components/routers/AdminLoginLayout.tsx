import React from "react";
import { Outlet, Navigate } from "react-router-dom";
import { useAuth } from "../contexts/AuthContext";

export default function AdminLoginLayout() {
  const { state } = useAuth();

  return (state.isAuthorized && state.isAdmin) ? <Navigate to="/admin" /> : <Outlet />;
}
