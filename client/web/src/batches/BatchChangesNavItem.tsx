import React from 'react'

import { NavItem, NavLink } from '@sourcegraph/wildcard/src/components/NavBar'

import { BatchChangesIconNav } from './icons'

interface Props {
    // Nothing for now.
}

/**
 * An item in {@link GlobalNavbar} that links to the batch changes area.
 */
export const BatchChangesNavItem: React.FunctionComponent<Props> = () => (
    <NavItem icon={BatchChangesIconNav}>
        <NavLink to="/batch-changes">Batch Changes</NavLink>
    </NavItem>
)
