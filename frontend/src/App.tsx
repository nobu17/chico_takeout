import React from "react";
import { Route, Routes, BrowserRouter } from "react-router-dom";
import {
  createTheme,
  responsiveFontSizes,
  ThemeProvider,
} from "@mui/material/styles";
import { AdapterDateFns } from "@mui/x-date-pickers/AdapterDateFns";
import { LocalizationProvider } from "@mui/x-date-pickers";
import { AdminAuthProvider } from "./components/contexts/AdminAuthContext";

import AdminAuthLayout from "./components/routers/AdminAuthLayout";
import AdminLoginLayout from "./components/routers/AdminLoginLayout";

import Container from "@mui/material/Container";
import Header from "./components/layouts/Header";

import Home from "./pages/Home";

import AdminHome from "./pages/admin/Home";
import AdminLogin from "./pages/admin/Login";
import AdminLogout from "./pages/admin/Logout";
import ItemKind from "./pages/admin/items/ItemKind";
import StockItem from "./pages/admin/items/StockItem";
import StockItemRemain from "./pages/admin/items/StockItemRemain";
import FoodItem from "./pages/admin/items/FoodItem";
import BusinessHour from "./pages/admin/stores/BusinessHour";
import SpecialBusinessHour from "./pages/admin/stores/SpecialBusinessHour";
import SpecialHoliday from "./pages/admin/stores/SpecialHoliday";
import "./App.css";

let theme = createTheme();
theme = responsiveFontSizes(theme);

function App() {
  return (
    <LocalizationProvider dateAdapter={AdapterDateFns}>
      <ThemeProvider theme={theme}>
        <AdminAuthProvider>
          <BrowserRouter>
            <Header />
            <Container maxWidth="lg">
              <Routes>
                <Route element={<AdminLoginLayout />}>
                  <Route path="/admin/sign_in" element={<AdminLogin />} />
                </Route>
                <Route element={<AdminAuthLayout />}>
                  <Route path="/admin" element={<AdminHome />} />
                  <Route path="/admin/sign_out" element={<AdminLogout />} />
                  <Route path="/admin/items/kind" element={<ItemKind />} />
                  <Route path="/admin/items/stock" element={<StockItem />} />
                  <Route
                    path="/admin/items/stock/remain"
                    element={<StockItemRemain />}
                  />
                  <Route
                    path="/admin/items/stock/remain"
                    element={<StockItemRemain />}
                  />
                  <Route path="/admin/items/food" element={<FoodItem />} />
                  <Route path="/admin/store/hour" element={<BusinessHour />} />
                  <Route
                    path="/admin/store/special_hour"
                    element={<SpecialBusinessHour />}
                  />
                  <Route
                    path="/admin/store/holiday"
                    element={<SpecialHoliday />}
                  />
                </Route>

                {/* <AdminLoginRoute
                  path="/admin/login"
                  element={<AdminLogin />}
                ></AdminLoginRoute>
                <AdminAuthRoute path="/admin" element={<AdminHome />} />
                <AdminAuthRoute
                  path="/admin/items/kind"
                  element={<ItemKind />}
                />
                <AdminAuthRoute
                  path="/admin/items/stock"
                  element={<StockItem />}
                />
                <AdminAuthRoute
                  path="/admin/items/stock/remain"
                  element={<StockItemRemain />}
                />
                <AdminAuthRoute
                  path="/admin/items/food"
                  element={<FoodItem />}
                />
                <AdminAuthRoute
                  path="/admin/store/hour"
                  element={<BusinessHour />}
                />
                <AdminAuthRoute
                  path="/admin/store/special_hour"
                  element={<SpecialBusinessHour />}
                />
                <AdminAuthRoute
                  path="/admin/store/holiday"
                  element={<SpecialHoliday />}
                /> */}
                <Route path="*" element={<Home />} />
              </Routes>
            </Container>
          </BrowserRouter>
        </AdminAuthProvider>
      </ThemeProvider>
    </LocalizationProvider>
  );
}

export default App;
