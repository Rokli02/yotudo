import { SvgIconComponent } from "@mui/icons-material";
import { MouseEventHandler } from "react";

export interface NavbarProps {
  items?: NavbarItem[];
}

export interface NavbarItem {
  navigateTo?: string;
  label: string;
  onClick?: MouseEventHandler<HTMLAnchorElement>
  Icon?: SvgIconComponent;
}