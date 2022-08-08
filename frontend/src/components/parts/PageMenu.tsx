import { NavLink as RouterLink } from "react-router-dom";

import ListSubheader from "@mui/material/ListSubheader";
import List from "@mui/material/List";
import ListItemButton from "@mui/material/ListItemButton";
import ListItemIcon from "@mui/material/ListItemIcon";
import ListItemText from "@mui/material/ListItemText";
import Stack from "@mui/material/Stack";
import DoubleArrowIcon from "@mui/icons-material/DoubleArrow";
import CoffeeIcon from "@mui/icons-material/Coffee";
import { Typography } from "@mui/material";

export type PageMenuProps = {
  title: string;
  icon: string;
  pageInfos: PageInfo[];
};

export type PageInfo = {
  title: string;
  url: string;
  disabled?: boolean;
};

const getIcon = (input: string) => {
  switch (input) {
    case "coffee": {
      return <CoffeeIcon />;
    }
    default: {
      return <CoffeeIcon />;
    }
  }
};

export default function PageMenu(props: PageMenuProps) {
  return (
    <>
      <List
        sx={{
          width: "100%",
          maxWidth: 460,
        }}
        component="nav"
        subheader={
          <ListSubheader
            component="div"
            sx={{
              color: "text.primary",
              fontSize: 20,
              textAlign: "center",
              fontWeight: "bold",
            }}
          >
            <Stack
              direction="row"
              justifyContent="center"
              alignItems="center"
              gap={1}
            >
              {getIcon(props.icon)}
              <Typography variant="h5">{props.title}</Typography>
            </Stack>
          </ListSubheader>
        }
      >
        {props.pageInfos.map((item, index) =>
          item.disabled ? (
            <ListItemButton key={index} disabled={true}>
              <ListItemIcon>
                <DoubleArrowIcon />
              </ListItemIcon>
              <ListItemText primary={item.title} />
            </ListItemButton>
          ) : (
            <ListItemButton key={index} component={RouterLink} to={item.url}>
              <ListItemIcon>
                <DoubleArrowIcon />
              </ListItemIcon>
              <ListItemText primary={item.title} />
            </ListItemButton>
          )
        )}
      </List>
    </>
  );
}
