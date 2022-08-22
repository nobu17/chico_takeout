import * as React from "react";
import { Link } from "react-router-dom";
import AppBar from "@mui/material/AppBar";
import Box from "@mui/material/Box";
import Toolbar from "@mui/material/Toolbar";
import IconButton from "@mui/material/IconButton";
import Typography from "@mui/material/Typography";
import Menu from "@mui/material/Menu";
import MenuIcon from "@mui/icons-material/Menu";
import Container from "@mui/material/Container";
import Button from "@mui/material/Button";
import MenuItem from "@mui/material/MenuItem";
import AdbIcon from "@mui/icons-material/Adb";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../contexts/AuthContext";

const title = "CHICO★SPICE";
const pages: NavItem[] = [
  { label: "テイクアウト予約", link: "/reserve" },
  { label: "マイページ", link: "/my_page", isUser: true },
  { label: "管理ページ", link: "/admin", isAdmin: true },
  { label: "ログアウト", link: "/auth/sign_out", isUser: true },
  { label: "店舗公式", link: "https://chico-sp-website.web.app" },
  { label: "イートイン予約", link: "https://nobu17.pythonanywhere.com" },
  { label: "お問い合わせ", link: "/inquiry" },
];

type NavItem = {
  label: string;
  link: string;
  isAdmin?: boolean;
  isUser?: boolean;
};

const Header = () => {
  const navigate = useNavigate();
  const { state } = useAuth();
  const [anchorElNav, setAnchorElNav] = React.useState(null);

  const getFilteredNaviItem = (): NavItem[] => {
    const menus = Array<NavItem>();

    for (const item of pages) {
      if (item.isAdmin) {
        if (state.isAdmin && state.isAuthorized) {
          menus.push(item);
        }
        continue;
      }
      if (item.isUser) {
        if (state.isAuthorized) {
          menus.push(item);
        }
        continue;
      }
      menus.push(item);
    }
    return menus;
  };

  const handleOpenNavMenu = (event: any) => {
    setAnchorElNav(event.currentTarget);
  };

  const handleCloseNavMenu = () => {
    setAnchorElNav(null);
  };

  const handleLinkClick = (url: string) => {
    setAnchorElNav(null);
    if (url.startsWith("http")) {
      window.open(url);
    } else {
      navigate(url);
    }
  };

  return (
    <AppBar position="static" sx={{ mb: 2 }}>
      <Container maxWidth="xl">
        <Toolbar disableGutters>
          <AdbIcon sx={{ display: { xs: "none", md: "flex" }, mr: 1 }} />
          <Typography
            variant="h6"
            noWrap
            component={Link}
            to="/"
            sx={{
              mr: 2,
              display: { xs: "none", md: "flex" },
              fontFamily: "monospace",
              fontWeight: 700,
              letterSpacing: ".3rem",
              color: "inherit",
              textDecoration: "none",
            }}
          >
            {title}
          </Typography>

          <Box sx={{ flexGrow: 1, display: { xs: "flex", md: "none" } }}>
            <IconButton
              size="large"
              aria-label="account of current user"
              aria-controls="menu-appbar"
              aria-haspopup="true"
              onClick={handleOpenNavMenu}
              color="inherit"
            >
              <MenuIcon />
            </IconButton>
            <Menu
              id="menu-appbar"
              anchorEl={anchorElNav}
              anchorOrigin={{
                vertical: "bottom",
                horizontal: "left",
              }}
              keepMounted
              transformOrigin={{
                vertical: "top",
                horizontal: "left",
              }}
              open={Boolean(anchorElNav)}
              onClose={handleCloseNavMenu}
              sx={{
                display: { xs: "block", md: "none" },
              }}
            >
              {getFilteredNaviItem().map((page) => (
                <MenuItem
                  key={page.label}
                  onClick={() => handleLinkClick(page.link)}
                >
                  <Typography textAlign="center">{page.label}</Typography>
                </MenuItem>
              ))}
            </Menu>
          </Box>
          <AdbIcon sx={{ display: { xs: "flex", md: "none" }, mr: 1 }} />
          <Typography
            variant="h5"
            noWrap
            component={Link}
            to="/"
            sx={{
              mr: 2,
              display: { xs: "flex", md: "none" },
              flexGrow: 1,
              fontFamily: "monospace",
              fontWeight: 700,
              letterSpacing: ".3rem",
              color: "inherit",
              textDecoration: "none",
            }}
          >
            {title}
          </Typography>
          <Box sx={{ flexGrow: 1, display: { xs: "none", md: "flex" } }}>
            {getFilteredNaviItem().map((page) => (
              <Button
                key={page.label}
                onClick={() => handleLinkClick(page.link)}
                sx={{ my: 2, color: "white", display: "block" }}
              >
                {page.label}
              </Button>
            ))}
          </Box>
        </Toolbar>
      </Container>
    </AppBar>
  );
};
export default Header;
