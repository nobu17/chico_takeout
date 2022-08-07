import React from "react";
import { Route, Routes, BrowserRouter } from "react-router-dom";
import {
  createTheme,
  responsiveFontSizes,
  ThemeProvider,
} from "@mui/material/styles";
import { AdapterDateFns } from "@mui/x-date-pickers/AdapterDateFns";
import ja from "date-fns/locale/ja";
import { LocalizationProvider } from "@mui/x-date-pickers";
import { AdminAuthProvider } from "./components/contexts/AuthContext";

import AdminAuthLayout from "./components/routers/AdminAuthLayout";
import AdminLoginLayout from "./components/routers/AdminLoginLayout";
import UserAuthLayout from "./components/routers/UserAuthLayout";
import UserLoginLayout from "./components/routers/UserLoginLayout";

import Container from "@mui/material/Container";
import Header from "./components/layouts/Header";

import Home from "./pages/Home";
import UserLogin from "./pages/auth/Login";
import Logout from "./pages/auth/Logout";

import MyHome from "./pages/mypage/Home";
import ReserveHistory from "./pages/mypage/reserve/ReserveHistory";

import ReserveHome from "./pages/reserve/Home";

import AdminHome from "./pages/admin/Home";
import AdminLogin from "./pages/admin/Login";

import ItemKind from "./pages/admin/items/ItemKind";
import StockItem from "./pages/admin/items/StockItem";
import StockItemRemain from "./pages/admin/items/StockItemRemain";
import FoodItem from "./pages/admin/items/FoodItem";
import BusinessHour from "./pages/admin/stores/BusinessHour";
import SpecialBusinessHour from "./pages/admin/stores/SpecialBusinessHour";
import SpecialHoliday from "./pages/admin/stores/SpecialHoliday";
import Calendar from "./pages/admin/orders/Calendar";

import "./App.css";

let theme = createTheme();
theme = responsiveFontSizes(theme);

function App() {
  return (
    <LocalizationProvider dateAdapter={AdapterDateFns} adapterLocale={ja}>
      <ThemeProvider theme={theme}>
        <AdminAuthProvider>
          <BrowserRouter>
            <Header />
            <Container maxWidth="lg">
              <Routes>
                <Route element={<UserLoginLayout />}>
                  <Route path="/auth/sign_in" element={<UserLogin />} />
                </Route>
                <Route element={<UserAuthLayout />}>
                  <Route path="/my_page" element={<MyHome />} />
                  <Route path="/my_page/history" element={<ReserveHistory />} />
                  <Route path="/auth/sign_out" element={<Logout />} />
                </Route>
                <Route element={<AdminLoginLayout />}>
                  <Route path="/admin/sign_in" element={<AdminLogin />} />
                </Route>
                <Route element={<AdminAuthLayout />}>
                  <Route path="/admin" element={<AdminHome />} />
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
                  <Route path="/admin/orders/calendar" element={<Calendar />} />
                </Route>
                <Route path="/reserve" element={<ReserveHome />} />
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
