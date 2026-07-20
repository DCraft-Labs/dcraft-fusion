import { defineConfig } from "vitepress";

export default defineConfig({
  title: "DCraft Fusion",
  description:
    "AI-assisted, engine-agnostic control plane for modern data platforms. One brain. Many muscles.",
  cleanUrls: true,
  ignoreDeadLinks: true,
  themeConfig: {
    logo: "/logo.svg",
    siteTitle: "DCraft Fusion",
    nav: [
      { text: "Guide", link: "/getting-started" },
      { text: "Architecture", link: "/architecture" },
      { text: "Open Core", link: "/open-core" },
      {
        text: "GitHub",
        link: "https://github.com/DCraft-Labs/dcraft-fusion",
      },
    ],
    sidebar: [
      {
        text: "Introduction",
        items: [
          { text: "What is Fusion", link: "/" },
          { text: "Getting started", link: "/getting-started" },
          { text: "Open core", link: "/open-core" },
        ],
      },
      {
        text: "Install",
        items: [
          { text: "Docker Compose", link: "/install/docker" },
          { text: "Helm", link: "/install/helm" },
        ],
      },
      {
        text: "Concepts",
        items: [
          { text: "Architecture", link: "/architecture" },
          { text: "CDC", link: "/cdc" },
          { text: "Authentication", link: "/authentication" },
        ],
      },
      {
        text: "Community",
        items: [
          { text: "Contributing", link: "/contributing" },
          { text: "Security", link: "/security" },
        ],
      },
    ],
    socialLinks: [
      {
        icon: "github",
        link: "https://github.com/DCraft-Labs/dcraft-fusion",
      },
    ],
    footer: {
      message: "Released under the Apache License 2.0.",
      copyright: "Copyright © DCraft Labs",
    },
    editLink: {
      pattern:
        "https://github.com/DCraft-Labs/dcraft-fusion/edit/main/apps/docs/:path",
      text: "Edit this page on GitHub",
    },
  },
});
