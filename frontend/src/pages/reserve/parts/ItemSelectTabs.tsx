import * as React from "react";
import { Grid, Box, Tabs, Tab } from "@mui/material";
import ContentPasteIcon from '@mui/icons-material/ContentPaste';
import ItemCard from "./ItemCard";
import Typography from "@mui/material/Typography";
import {
  Cart,
  ItemInfo,
  ItemRequest,
  OptionItemInfo,
} from "../../../hooks/UseItemCart";

type ItemSelectTabsProps = {
  allItems: CategoryItems[];
  cart: Cart;
  onRequestChanged?: callback;
};
interface callback {
  (item: ItemRequest): void;
}

type CategoryItems = {
  title: string;
  items: ItemInfo[];
};

export default function ItemSelectTabs(props: ItemSelectTabsProps) {
  const getQuantity = (itemId: string): number => {
    if (props.cart.items[itemId]) {
      return props.cart.items[itemId].quantity;
    }
    return 0;
  };
  const getSelectedOptions = (itemId: string): OptionItemInfo[] => {
    if (props.cart.items[itemId]) {
      return props.cart.items[itemId].selectOptions;
    }
    return [];
  };
  const [selectedTab, setSelectedTab] = React.useState(0);

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setSelectedTab(newValue);
  };

  return (
    <>
      <Box sx={{ maxWidth: { xs: 320, sm: 1000 }, bgcolor: "background.paper" }}>
        <Tabs
          value={selectedTab}
          onChange={handleTabChange}
          variant="scrollable"
          scrollButtons
          allowScrollButtonsMobile
          aria-label="item"
        >
          {props.allItems.map((items, index) => {
            return <Tab label={items.title} key={index} icon={<ContentPasteIcon />} />;
          })}
        </Tabs>
      </Box>
      {props.allItems.map((items, index) => {
        return (
          <div key={index} role="tabpanel" hidden={selectedTab !== index}>
            {selectedTab === index && (
              <>
                <Typography
                  component="h3"
                  variant="h4"
                  align="center"
                  color="text.primary"
                  gutterBottom
                  sx={{
                    mt: 3,
                  }}
                >
                  {items.title}
                </Typography>
                <Grid container spacing={2}>
                  {items.items.map((item, index) => {
                    return (
                      <Grid item xs={12} md={6} key={index + "item"}>
                        <ItemCard
                          item={item}
                          quantity={getQuantity(item.id)}
                          selectedOptions={getSelectedOptions(item.id)}
                          onChanged={props.onRequestChanged}
                        ></ItemCard>
                      </Grid>
                    );
                  })}
                </Grid>
              </>
            )}
          </div>
        );
      })}
    </>
  );
}
