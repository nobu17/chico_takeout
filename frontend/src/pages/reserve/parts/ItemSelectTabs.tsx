import * as React from "react";
import { Grid, Box, Tabs, Tab } from "@mui/material";
import ContentPasteIcon from "@mui/icons-material/ContentPaste";
import ItemSelectCard from "./ItemSelectCard";
import {
  Cart,
  ItemInfo,
  ItemRequest,
  OptionItemInfo,
} from "../../../hooks/UseItemCart";
import CartButton from "./CartButton";

type ItemSelectTabsProps = {
  allItems: CategoryItems[];
  cart: Cart;
  onRequestChanged: callback;
  onCartUpdated: (cart: Cart) => void;
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

  const handleCartUpdated = (cart: Cart) => {
    props.onCartUpdated(cart);
  };

  return (
    <>
      <Box
        sx={{
          maxWidth: { xs: 320, sm: 1000 },
          position: "sticky",
          top: 0,
          bgcolor: "background.paper",
        }}
      >
        <Tabs
          value={selectedTab}
          onChange={handleTabChange}
          variant="scrollable"
          scrollButtons
          allowScrollButtonsMobile
          aria-label="item"
        >
          {props.allItems.map((items, index) => {
            return (
              <Tab
                label={items.title}
                key={index}
                icon={<ContentPasteIcon />}
              />
            );
          })}
        </Tabs>
      </Box>
      {props.allItems.map((items, index) => {
        return (
          <div key={index} role="tabpanel" hidden={selectedTab !== index}>
            {selectedTab === index && (
              <>
                <Grid container spacing={2} alignItems="stretch">
                  {items.items.map((item, index) => {
                    return (
                      <Grid item xs={6} md={4} key={index + "item"}>
                        <ItemSelectCard
                          item={item}
                          quantity={getQuantity(item.id)}
                          selectedOptions={getSelectedOptions(item.id)}
                          onChanged={props.onRequestChanged}
                        ></ItemSelectCard>
                      </Grid>
                    );
                  })}
                </Grid>
              </>
            )}
          </div>
        );
      })}
      <CartButton
        allItems={props.allItems}
        cart={props.cart}
        onUpdated={handleCartUpdated}
      />
    </>
  );
}
