function overflowMenu(tree = document) {
    tree.querySelectorAll("[data-overflow-menu]").forEach(menuRoot => {
      const
      button = menuRoot.querySelector("[aria-haspopup]"),
      menu = menuRoot.querySelector("[role=menu]"),
      items = [...menu.querySelectorAll("[role=menuitem]")];
    });
}
  
addEventListener("htmx:load", e => overflowMenu(e.target));