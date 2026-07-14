# Visual Acceptance Test Screenshots

This directory contains screenshots captured during the KDSE implementation acceptance testing for the Forge application.

## Screenshots

### 01-welcome-screen.png
**Title:** Welcome Screen  
**Description:** The initial landing page of the Forge application showing the welcome screen with options to "Load Utility-Scale Solar Farm" or "Open Existing Project".  
**Functionality Demonstrated:** Application bootstrap and initial user interface.

### 02-reference-world-loaded.png
**Title:** Reference World Loaded  
**Description:** The Utility-Scale Solar Farm reference world has been loaded, showing the Plant Explorer with the full equipment hierarchy (Environment, Grid Connection, Substation, Switchyard, PV Arrays, Revenue Meters, Station Loads).  
**Functionality Demonstrated:** World loading and plant explorer navigation structure.

### 03-operation-workspace.png
**Title:** Operation Workspace  
**Description:** Full view of the Operation Workspace showing the Plant Explorer, Single Line Diagram with live simulation data, and the Equipment Details panel. The simulation is in stopped state.  
**Functionality Demonstrated:** Main operation workspace layout with simulation controls.

### 04-single-line-diagram.png
**Title:** Single Line Diagram  
**Description:** Close view of the Single Line Diagram showing the solar farm topology with Grid Connection, Revenue Meter (0 kW), Transformer (0%), Bus (480 V), PV Arrays (2.1 values), and Station Loads (50 kW).  
**Functionality Demonstrated:** Live simulation data visualization in the single line diagram.

### 05-equipment-selected.png
**Title:** Equipment Selected  
**Description:** PV Array 1 (2.5 MW) has been selected from the plant explorer, showing the Equipment Details panel with the Explain tab active. Displays current power generation (2.1 kW), solar irradiance (850 W/m²), array capacity (2500 kW), and panel efficiency (19%).  
**Functionality Demonstrated:** Equipment selection and detail view with live simulation data.

### 06-equipment-details-explain.png
**Title:** Equipment Details - Explain Tab  
**Description:** The Equipment Details panel showing the Explain tab for PV Array 1. Includes explanation of power generation factors and how the simulation influences the readings.  
**Functionality Demonstrated:** AI-generated explanations for equipment behavior based on simulation state.

### 07-analysis-panel.png
**Title:** Analysis Panel  
**Description:** The Analysis tab view showing Timeline, Events, and Why? sections. Displays 0 Active Alarms and 1 Total Event (Simulation Started).  
**Functionality Demonstrated:** Analysis and monitoring capabilities of the simulation.

### 08-simulation-running.png
**Title:** Simulation Running  
**Description:** The simulation is actively running (status shows "Running" and time has progressed to 00:00:08). The PV arrays show 0.0 kW because the sun is below the horizon at the current simulation time.  
**Functionality Demonstrated:** Simulation execution and live data updates based on time-of-day lighting conditions.

### 09-simulation-speed-8x.png
**Title:** Simulation Speed Control  
**Description:** The simulation speed selector showing available speeds (0.1x, 0.25x, 0.5x, 1x, 2x, 4x, 8x). The simulation continues to run with time progressing rapidly.  
**Functionality Demonstrated:** Simulation speed control for fast-forwarding through time.

### 10-final-overview.png
**Title:** Final Overview  
**Description:** Complete overview of the Operation Workspace with all panels visible, the plant hierarchy expanded, and the simulation running. Shows the complete KDSE data flow from simulation to UI.  
**Functionality Demonstrated:** Full application overview demonstrating the complete KDSE implementation.

---

## Testing Notes

### Build Status
✅ UI build completed successfully (`npm run build`)

### Runtime Status
✅ Application runs without JavaScript errors  
✅ No console errors detected

### Functional Verification
- **Welcome Screen:** Renders correctly
- **World Loading:** Reference world loads without errors
- **Plant Explorer:** Tree navigation works correctly
- **Single Line Diagram:** Displays live simulation data
- **Equipment Details:** Shows real-time measurements from simulation
- **Analysis Panel:** Displays timeline events and statistics
- **Simulation Controls:** Play, pause, and speed controls work
- **Simulation State:** Values update based on sun position and irradiance

### KDSE Data Flow Verification
✅ Simulation → Measurements → Operation Workspace → Equipment Details → Analysis → Single Line Diagram

All engineering measurements visible in Operation Mode originate from the simulation. The simulation remains the single source of truth.
