<script lang="ts">
  import { page } from '$app/stores'
  import './layout.css'

  let { children } = $props()
   $inspect($page.data) 
</script>

{@render children()}
