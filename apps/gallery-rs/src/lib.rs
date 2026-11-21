//! Browser client of the Rust gallery.
//!
//! For now it only proves the build pipeline: it finds the `#stage`
//! canvas and paints it a single colour.

use wasm_bindgen::prelude::*;
use web_sys::{CanvasRenderingContext2d, HtmlCanvasElement};

/// Runs once, as soon as the module is instantiated.
#[wasm_bindgen(start)]
fn start() -> Result<(), JsValue> {
    let document = web_sys::window()
        .and_then(|w| w.document())
        .ok_or("no document")?;
    let canvas: HtmlCanvasElement = document
        .get_element_by_id("stage")
        .ok_or("no canvas with id \"stage\"")?
        .dyn_into()?;
    let ctx: CanvasRenderingContext2d = canvas
        .get_context("2d")?
        .ok_or("canvas has no 2d context")?
        .dyn_into()?;

    ctx.set_fill_style_str("#6d4aff");
    ctx.fill_rect(0.0, 0.0, canvas.width().into(), canvas.height().into());
    Ok(())
}
